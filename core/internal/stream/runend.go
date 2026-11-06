package stream

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wandb/wandb/core/internal/filestream"
	"github.com/wandb/wandb/core/internal/observability"
)

// RunEndFileName is the name of the run-end record in the run's sync
// directory.
const RunEndFileName = "run_end.json"

// End reasons written to the run-end record.
const (
	RunEndFinished      = "finished"
	RunEndUploadStopped = "upload_stopped"
	RunEndShutdown      = "shutdown"
)

// runEndRecord is the JSON content of the run-end record.
//
// It is a local, best-effort account of how a run ended and how much of its
// data the server acknowledged, so that a later reader (a person, `wandb sync`
// or another tool) can tell a complete run from an incomplete one.
type runEndRecord struct {
	Version int `json:"version"`

	// EndReason is RunEndFinished, RunEndUploadStopped or RunEndShutdown.
	EndReason string `json:"end_reason"`

	// ShutdownReason says what started a service shutdown, such as a signal
	// or the parent process exiting. Empty if there was no shutdown.
	ShutdownReason string `json:"shutdown_reason,omitempty"`

	// ShutdownDrained is whether a shutdown finished the run before its
	// grace period expired.
	ShutdownDrained *bool `json:"shutdown_drained,omitempty"`

	// ExitCode is the run's exit code, if it got as far as exiting.
	ExitCode *int32 `json:"exit_code,omitempty"`

	// UploadError is the error that stopped filestream uploads, if any.
	UploadError string `json:"upload_error,omitempty"`

	// LostRecords counts records that could not be read back from the
	// transaction log and so were never uploaded.
	LostRecords int64 `json:"lost_records,omitempty"`

	// Offline is whether the run never uploaded data while it ran.
	Offline bool `json:"offline,omitempty"`

	// CompleteAcked is whether the server acknowledged the request that
	// marked the run complete.
	CompleteAcked bool `json:"complete_acked"`

	// AckedOffsets is, per filestream file, the offset just past the last
	// line the server acknowledged.
	AckedOffsets map[string]int `json:"acked_offsets,omitempty"`

	UpdatedAt string `json:"updated_at"`
}

// RunEndRecorder maintains a run's run-end record on disk.
//
// Every method rewrites the whole file atomically, so the file always holds
// the latest known state even if the process is killed right after.
//
// A nil *RunEndRecorder is valid and records nothing.
type RunEndRecorder struct {
	mu     sync.Mutex
	path   string
	logger *observability.CoreLogger
	record runEndRecord

	shutdownStarted bool
	exited          bool
}

// NewRunEndRecorder returns a recorder that writes into syncDir.
//
// Returns nil if syncDir is empty.
func NewRunEndRecorder(
	syncDir string,
	logger *observability.CoreLogger,
) *RunEndRecorder {
	if syncDir == "" {
		return nil
	}

	return &RunEndRecorder{
		path:   filepath.Join(syncDir, RunEndFileName),
		logger: logger,
		record: runEndRecord{Version: 1},
	}
}

// Path returns the path of the run-end record.
func (r *RunEndRecorder) Path() string {
	if r == nil {
		return ""
	}
	return r.path
}

// UploadStopped records that filestream stopped uploading because of err.
func (r *RunEndRecorder) UploadStopped(err error, progress filestream.Progress) {
	if r == nil || err == nil {
		return
	}

	r.update(func(rec *runEndRecord) {
		if rec.UploadError == "" {
			rec.UploadError = err.Error()
		}
		r.applyProgress(rec, progress)
	})
}

// RecordsLost records that n saved records were skipped and never uploaded.
func (r *RunEndRecorder) RecordsLost(n int64) {
	if r == nil || n <= 0 {
		return
	}

	r.update(func(rec *runEndRecord) {
		rec.LostRecords += n
	})
}

// Exited records the run's exit and its final upload state.
//
// fsErr and progress describe the filestream; offline is true when there
// was no filestream.
func (r *RunEndRecorder) Exited(
	exitCode int32,
	offline bool,
	fsErr error,
	progress filestream.Progress,
) {
	if r == nil {
		return
	}

	r.update(func(rec *runEndRecord) {
		r.exited = true
		rec.ExitCode = &exitCode
		rec.Offline = offline
		if fsErr != nil && rec.UploadError == "" {
			rec.UploadError = fsErr.Error()
		}
		r.applyProgress(rec, progress)
	})
}

// ShutdownStarted records that a service shutdown began for reason.
func (r *RunEndRecorder) ShutdownStarted(reason string) {
	if r == nil {
		return
	}

	r.update(func(rec *runEndRecord) {
		r.shutdownStarted = true
		rec.ShutdownReason = reason
		drained := false
		rec.ShutdownDrained = &drained
	})
}

// ShutdownFinished records whether the shutdown drained the run in time.
func (r *RunEndRecorder) ShutdownFinished(drained bool) {
	if r == nil {
		return
	}

	r.update(func(rec *runEndRecord) {
		if !r.shutdownStarted {
			return
		}
		rec.ShutdownDrained = &drained
	})
}

// applyProgress copies filestream progress into the record.
//
// Offsets only move forward, so a stale snapshot never hides acked data.
func (r *RunEndRecorder) applyProgress(
	rec *runEndRecord,
	progress filestream.Progress,
) {
	if progress.CompleteAcked {
		rec.CompleteAcked = true
	}

	if len(progress.AckedOffsets) == 0 {
		return
	}

	if rec.AckedOffsets == nil {
		rec.AckedOffsets = make(map[string]int, len(progress.AckedOffsets))
	}
	for name, offset := range progress.AckedOffsets {
		rec.AckedOffsets[name] = max(rec.AckedOffsets[name], offset)
	}
}

// endReason derives the record's end reason from what has been recorded.
func (r *RunEndRecorder) endReason(rec *runEndRecord) string {
	switch {
	case r.shutdownStarted:
		return RunEndShutdown
	case rec.UploadError != "" || rec.LostRecords > 0:
		return RunEndUploadStopped
	case r.exited:
		return RunEndFinished
	default:
		return ""
	}
}

// update applies change to the record and rewrites the file.
func (r *RunEndRecorder) update(change func(rec *runEndRecord)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	change(&r.record)
	r.record.EndReason = r.endReason(&r.record)
	r.record.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)

	snapshot := r.record
	snapshot.AckedOffsets = maps.Clone(r.record.AckedOffsets)

	if err := writeFileAtomically(r.path, snapshot); err != nil && r.logger != nil {
		r.logger.Error("runend: failed to write run-end record", "error", err)
	}
}

// writeFileAtomically writes rec as JSON to path through a synced temporary
// file, so readers never see a partial record.
func writeFileAtomically(path string, rec runEndRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("runend: marshal: %v", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), RunEndFileName+".tmp*")
	if err != nil {
		return fmt.Errorf("runend: create temp file: %v", err)
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	if _, err := tmp.Write(append(data, '\n')); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("runend: write: %v", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("runend: sync: %v", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("runend: close: %v", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("runend: rename: %v", err)
	}
	return nil
}
