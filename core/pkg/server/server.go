package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/wandb/wandb/core/internal/monitor"
	"github.com/wandb/wandb/core/internal/runsync"
	"github.com/wandb/wandb/core/internal/stream"
	"github.com/wandb/wandb/core/internal/sweeps/scheduler"
	"github.com/wandb/wandb/core/pkg/server/listeners"
)

const (
	BufferSize                         = 32
	IntervalCheckParentPidMilliseconds = 100
)

// Server is the top-level object for the wandb-core process.
type Server struct {
	// connNumber is the number of the last connection created.
	//
	// It is used to generate unique IDs for connections.
	connNumber atomic.Int64

	// serverLifetimeCtx is cancelled when the server should shut down.
	serverLifetimeCtx context.Context

	// stopServer cancels serverLifetimeCtx.
	stopServer context.CancelFunc

	// forceShutdown reports whether Serve should return without waiting for
	// in-flight connections to finish.
	forceShutdown       atomic.Bool
	gracefulShutdown    atomic.Bool
	shutdownOnce        sync.Once
	shutdownDone        chan struct{}
	forceDone           chan struct{}
	forceOnce           sync.Once
	shutdownGracePeriod time.Duration

	// streamMux maps stream IDs to streams.
	streamMux *stream.StreamMux

	// runSyncManager implements `wandb sync` operations.
	runSyncManager *runsync.RunSyncManager

	// sweepSchedBroker implements `wandb sweep-scheduler` sessions.
	sweepSchedBroker *scheduler.IPCSessionBroker

	// symonResourceManager manages the system metrics sidecar.
	symonResourceManager *monitor.SymonResourceManager

	// connectionsWG is the WaitGroup to wait for all connections to finish
	// and for the serve goroutine to finish
	connectionsWG sync.WaitGroup

	// parentPID is parent process's ID.
	//
	// The server exits if the parent process is gone.
	parentPID int

	// detached is whether the server should ignore parent process lifetime.
	detached bool

	// idleTimeout specifies how long the server should stay alive without any
	// connected clients before shutting down in detached mode.
	idleTimeout time.Duration

	// activeConnections is the number of currently connected clients.
	activeConnections atomic.Int64

	// idleTimer is started when the server has zero connected clients.
	idleTimer *time.Timer

	// idleTimerMu protects idleTimer.
	idleTimerMu sync.Mutex

	// idleShutdownStarted reports whether the idle timer callback has started
	// shutting the server down.
	idleShutdownStarted bool

	// commit is the W&B Git commit hash.
	commit string

	// listenOnLocalhost is whether to open a localhost socket even if Unix
	// sockets are supported.
	listenOnLocalhost bool

	// loggerPath is the default logger path
	loggerPath string

	// logLevel is the log level
	logLevel slog.Level
}

// ForceStop shuts down the server without waiting for in-flight work.
func (s *Server) ForceStop() {
	s.forceShutdown.Store(true)
	s.forceOnce.Do(func() { close(s.forceDone) })
	s.stopServer()
}

// finalLogSyncTimeout bounds the last transaction-log fsync of a shutdown,
// which runs even after the grace period expires or a second signal arrives.
const finalLogSyncTimeout = time.Second

// Shutdown gives active runs a bounded opportunity to flush. It stops IPC
// admission first; a second signal can bypass the grace period via ForceStop.
//
// The reason is written into each run's run-end record.
func (s *Server) Shutdown(reason string) {
	s.shutdownOnce.Do(func() {
		s.gracefulShutdown.Store(true)

		// Take the streams now: finishing them removes them from the mux.
		streams := s.streamMux.Streams()
		for _, st := range streams {
			st.RecordShutdownStarted(reason)
		}

		s.stopServer()
		go func() {
			defer close(s.shutdownDone)
			if s.shutdownGracePeriod <= 0 {
				s.ForceStop()
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), s.shutdownGracePeriod)
			defer cancel()
			done := make(chan struct{})
			go func() {
				defer close(done)
				// Save already-buffered records even if a connection is stuck on I/O.
				stream.SyncTransactionLogs(streams)
				select {
				case <-s.waitForConnectionsToFinish():
				case <-ctx.Done():
					return
				}
				s.streamMux.FinishAndCloseAllStreams(255)
			}()

			drained := false
			select {
			case <-done:
				drained = true
				slog.Info("server: shutdown drain completed")
			case <-ctx.Done():
				slog.Warn(
					"server: shutdown grace period expired; local data may need wandb sync --include-online",
				)
			case <-s.forceDone:
				slog.Warn("server: shutdown forced; local data may need wandb sync --include-online")
			}

			// Save records that arrived while draining. A finished stream's
			// log is already closed and synced, so this is a no-op for it.
			syncLogsWithin(streams, finalLogSyncTimeout)
			for _, st := range streams {
				st.RecordShutdownFinished(drained)
			}
		}()
	})
}

// syncLogsWithin fsyncs the streams' transaction logs, giving up after
// timeout so that a stuck filesystem cannot block process exit.
func syncLogsWithin(streams []*stream.Stream, timeout time.Duration) {
	synced := make(chan struct{})
	go func() {
		defer close(synced)
		stream.SyncTransactionLogs(streams)
	}()

	select {
	case <-synced:
	case <-time.After(timeout):
		slog.Warn("server: timed out syncing transaction logs during shutdown")
	}
}

type ServerParams struct {
	Commit              string
	EnableDCGMProfiling bool
	ListenOnLocalhost   bool
	LoggerPath          string
	LogLevel            slog.Level
	ParentPID           int
	Detached            bool
	IdleTimeout         time.Duration
	ShutdownGracePeriod time.Duration
}

func NewServer(params ServerParams) *Server {
	serverLifetimeCtx, stopServer := context.WithCancel(context.Background())

	return &Server{
		serverLifetimeCtx:    serverLifetimeCtx,
		shutdownDone:         make(chan struct{}),
		forceDone:            make(chan struct{}),
		shutdownGracePeriod:  params.ShutdownGracePeriod,
		stopServer:           stopServer,
		streamMux:            stream.NewStreamMux(),
		runSyncManager:       runsync.NewRunSyncManager(),
		sweepSchedBroker:     scheduler.NewIPCSessionBroker(scheduler.NewTaskResolver),
		symonResourceManager: monitor.NewSymonResourceManager(params.EnableDCGMProfiling),
		connectionsWG:        sync.WaitGroup{},
		parentPID:            params.ParentPID,
		detached:             params.Detached,
		idleTimeout:          params.IdleTimeout,
		commit:               params.Commit,
		listenOnLocalhost:    params.ListenOnLocalhost,
		loggerPath:           params.LoggerPath,
		logLevel:             params.LogLevel,
	}
}

// exitWhenParentIsGone exits the process if the parent process is killed.
func (s *Server) exitWhenParentIsGone() {
	slog.Info("server: will exit if parent process dies", "ppid", s.parentPID)

	ticker := time.NewTicker(IntervalCheckParentPidMilliseconds * time.Millisecond)
	defer ticker.Stop()
	for os.Getppid() == s.parentPID {
		select {
		case <-s.serverLifetimeCtx.Done():
			return
		case <-ticker.C:
		}
	}

	slog.Info("server: parent process exited, terminating service process")

	s.Shutdown("parent_exited")
}

func closeAll(listenerList []net.Listener) {
	for idx, listener := range listenerList {
		if err := listener.Close(); err != nil {
			slog.Error("server: failed to close listener", "index", idx, "error", err)
		}
	}
}

func (s *Server) Serve(portFile string) error {
	listenerList, portInfo, err := listeners.Config{
		ParentPID:         s.parentPID,
		ListenOnLocalhost: s.listenOnLocalhost,
	}.MakeListeners()

	if err != nil {
		return err
	}

	var closeAllOnce sync.Once
	closeListeners := func() {
		closeAllOnce.Do(func() { closeAll(listenerList) })
	}
	defer closeListeners()

	if err := portInfo.WriteToFile(portFile); err != nil {
		return err
	}

	if s.parentPID != 0 && !s.detached {
		go s.exitWhenParentIsGone()
	}

	if s.detached && s.idleTimeout > 0 {
		// Start the idle timer immediately. If a client connects, it will be stopped.
		s.startIdleTimer()
	}

	for _, listener := range listenerList {
		s.connectionsWG.Go(func() {
			s.acceptConnections(listener)
		})
	}

	// Wait for the signal to shut down.
	<-s.serverLifetimeCtx.Done()
	slog.Info("server: is shutting down")

	s.stopIdleTimer()

	// Stop accepting new connections.
	closeListeners()

	// A graceful shutdown ends with a bounded log sync even if it is forced
	// part way, so wait for it before letting the process exit.
	if s.gracefulShutdown.Load() {
		<-s.shutdownDone
		return nil
	}

	if s.forceShutdown.Load() {
		slog.Info("server: forced shutdown")
		return nil
	}
	<-s.waitForConnectionsToFinish()
	slog.Info("server: all connections closed")
	return nil
}

func (s *Server) waitForConnectionsToFinish() chan struct{} {
	done := make(chan struct{})
	go func() {
		s.connectionsWG.Wait()
		close(done)
	}()

	return done
}

// acceptConnections accepts incoming connections on the listener.
//
// It blocks until the listener is closed or an error is encountered.
func (s *Server) acceptConnections(listener net.Listener) {
	slog.Info("server: accepting connections", "addr", listener.Addr())

	for {
		conn, err := listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			slog.Info("server: listener closed", "addr", listener.Addr())
			return
		}

		// If a connection is added to the listen queue but closed before we
		// reach here (like if the client suddenly shuts down), the accept()
		// system call returns ECONNABORTED.
		//
		// EMFILE, ENFILE mean we're out of file descriptors.
		// ENOBUFS, ENOMEM mean we're out of memory.
		// These conditions might be temporary, so we sleep and try again.
		if errors.Is(err, syscall.ECONNABORTED) ||
			errors.Is(err, syscall.EMFILE) ||
			errors.Is(err, syscall.ENFILE) ||
			errors.Is(err, syscall.ENOBUFS) ||
			errors.Is(err, syscall.ENOMEM) {
			slog.Warn(
				"server: failed to accept connection",
				"addr", listener.Addr(),
				"error", err,
			)
			time.Sleep(time.Second)
			continue
		}

		// Give up and shut down on any unexpected error.
		if err != nil {
			slog.Error("server: unknown error accepting connection", "error", err)
			s.stopServer()
			return
		}

		s.connectionsWG.Go(func() {
			s.onConnectionStart()
			defer s.onConnectionEnd()
			s.handleConnection(conn)
		})
	}
}

// handleConnection processes incoming data on a connection.
//
// This blocks until the connection closes.
func (s *Server) handleConnection(conn net.Conn) {
	var id string
	if addr := conn.RemoteAddr().String(); addr != "" {
		id = fmt.Sprintf("%d(%s)", s.connNumber.Add(1), addr)
	} else {
		id = fmt.Sprintf("%d", s.connNumber.Add(1))
	}

	NewConnection(
		s.serverLifetimeCtx,
		s.stopServer,
		ConnectionParams{
			ID:                   id,
			Conn:                 conn,
			StreamMux:            s.streamMux,
			RunSyncManager:       s.runSyncManager,
			SweepSchedBroker:     s.sweepSchedBroker,
			SymonResourceManager: s.symonResourceManager,
			Commit:               s.commit,
			LoggerPath:           s.loggerPath,
			LogLevel:             s.logLevel,
		},
	).ManageConnectionData()
}

func (s *Server) onConnectionStart() {
	if !s.detached || s.idleTimeout <= 0 {
		return
	}

	// Stop the idle timer when the first client connects.
	if s.activeConnections.Add(1) == 1 {
		s.stopIdleTimer()
	}
}

func (s *Server) onConnectionEnd() {
	if !s.detached || s.idleTimeout <= 0 {
		return
	}

	// Start the idle timer when the last client disconnects.
	if s.activeConnections.Add(-1) == 0 {
		s.startIdleTimer()
	}
}

func (s *Server) startIdleTimer() {
	s.idleTimerMu.Lock()
	defer s.idleTimerMu.Unlock()

	if s.idleShutdownStarted ||
		(s.idleTimer != nil && !s.idleTimer.Stop()) {
		return
	}

	s.idleTimer = time.AfterFunc(s.idleTimeout, func() {
		s.idleTimerMu.Lock()
		s.idleShutdownStarted = true
		s.idleTimer = nil
		s.idleTimerMu.Unlock()

		slog.Info("server: idle timeout reached, shutting down", "idle-timeout", s.idleTimeout)
		s.stopServer()
		// Gracefully finish and close all streams.
		s.streamMux.FinishAndCloseAllStreams(0)
	})
}

func (s *Server) stopIdleTimer() {
	s.idleTimerMu.Lock()
	defer s.idleTimerMu.Unlock()

	if s.idleShutdownStarted || s.idleTimer == nil {
		return
	}

	if !s.idleTimer.Stop() {
		s.idleShutdownStarted = true
		s.idleTimer = nil
		return
	}

	s.idleTimer = nil
}
