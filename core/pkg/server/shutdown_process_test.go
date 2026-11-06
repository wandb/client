//go:build !windows

package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/wandb/wandb/core/internal/observability"
	"github.com/wandb/wandb/core/internal/transactionlog"
)

// This exercises the Python SDK, IPC, service shutdown, and transaction log.
// Set WANDB_TEST_PYTHON and WANDB_TEST_CORE_BINARY to opt in.
func TestShutdownProcessPreservesOfflineHistory(t *testing.T) {
	python := os.Getenv("WANDB_TEST_PYTHON")
	if python == "" || os.Getenv("WANDB_TEST_CORE_BINARY") == "" {
		t.Skip("requires built core and Python SDK dependencies including psutil")
	}
	sdkRoot, err := filepath.Abs("../../..")
	require.NoError(t, err)
	scenarios := []string{"sigterm", "parent_death", "sigkill_control"}
	if runtime.GOOS == "linux" {
		scenarios = append(scenarios, "parent_death_prctl")
	}
	for _, scenario := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			logFile := runShutdownScenario(t, python, sdkRoot, scenario)
			requireLogContents(t, logFile, scenario)
			requireRunEnd(t, filepath.Dir(logFile), scenario)
		})
	}
}

// shutdownReady is what testdata/shutdown.py writes once it has logged.
type shutdownReady struct {
	PID  int    `json:"pid"`
	File string `json:"file"`
}

// runShutdownScenario runs testdata/shutdown.py, stops its wandb-core the way
// the scenario says, waits for wandb-core to exit and returns the path of the
// run's transaction log.
func runShutdownScenario(t *testing.T, python, sdkRoot, scenario string) string {
	t.Helper()
	root := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	worker := exec.CommandContext(ctx, python, "testdata/shutdown.py")
	worker.Env = append(
		os.Environ(),
		"PYTHONPATH="+sdkRoot,
		"M2_DIR="+root,
		"WANDB_SILENT=true",
	)
	if scenario == "parent_death_prctl" {
		worker.Env = append(worker.Env, "M2_OS_PID_SHUTDOWN=1")
	}
	log, err := os.Create(filepath.Join(root, "worker.log"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = log.Close() })
	worker.Stdout, worker.Stderr = log, log
	require.NoError(t, worker.Start())
	t.Cleanup(func() { _ = worker.Process.Kill(); _ = worker.Wait() })

	var ready shutdownReady
	require.Eventually(t, func() bool {
		data, err := os.ReadFile(filepath.Join(root, "ready.json"))
		return err == nil && json.Unmarshal(data, &ready) == nil
	}, 20*time.Second, 20*time.Millisecond)
	t.Cleanup(func() { _ = syscall.Kill(ready.PID, syscall.SIGKILL) })

	switch scenario {
	case "parent_death", "parent_death_prctl":
		require.NoError(t, worker.Process.Kill())
	case "sigkill_control":
		require.NoError(t, syscall.Kill(ready.PID, syscall.SIGKILL))
	default:
		require.NoError(t, syscall.Kill(ready.PID, syscall.SIGTERM))
	}

	// psutil handles zombies as well as processes already reaped by the parent.
	wait := exec.CommandContext(ctx, python, "-c", `import psutil,sys,time
p=int(sys.argv[1])
while psutil.pid_exists(p):
 try:
  if psutil.Process(p).status()==psutil.STATUS_ZOMBIE: break
 except psutil.NoSuchProcess: break
 time.sleep(.02)`, strconv.Itoa(ready.PID))
	require.NoError(t, wait.Run())
	return ready.File
}

// requireLogContents checks the history rows and exit record in the log.
func requireLogContents(t *testing.T, logFile, scenario string) {
	t.Helper()
	rows, exitCode := readHistoryAndExit(t, logFile, scenario == "sigkill_control")

	if scenario == "sigkill_control" {
		require.Less(t, rows, 25)
		require.Nil(t, exitCode)
		return
	}
	require.Equal(t, 25, rows)
	require.NotNil(t, exitCode)
	require.EqualValues(t, 255, *exitCode)
}

// readHistoryAndExit counts history records and returns the exit code, if any.
//
// If allowTruncated, a partially written last block ends the read.
func readHistoryAndExit(
	t *testing.T,
	logFile string,
	allowTruncated bool,
) (int, *int32) {
	t.Helper()
	reader, err := transactionlog.OpenReader(logFile, observability.NewNoOpLogger())
	require.NoError(t, err)
	defer reader.Close()

	var rows int
	var exitCode *int32
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) ||
			(allowTruncated && errors.Is(err, io.ErrUnexpectedEOF)) {
			return rows, exitCode
		}
		require.NoError(t, err)
		if record.GetHistory() != nil {
			rows++
		}
		if exit := record.GetExit(); exit != nil {
			value := exit.ExitCode
			exitCode = &value
		}
	}
}

// requireRunEnd checks the run-end record left by the shutdown.
func requireRunEnd(t *testing.T, runDir, scenario string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(runDir, "run_end.json"))
	if scenario == "sigkill_control" {
		require.ErrorIs(t, err, os.ErrNotExist, "SIGKILL runs no shutdown code")
		return
	}
	require.NoError(t, err)

	var runEnd struct {
		EndReason       string `json:"end_reason"`
		ShutdownReason  string `json:"shutdown_reason"`
		ShutdownDrained *bool  `json:"shutdown_drained"`
		ExitCode        *int32 `json:"exit_code"`
		Offline         bool   `json:"offline"`
	}
	require.NoError(t, json.Unmarshal(data, &runEnd))

	wantReason := "parent_exited"
	if scenario == "sigterm" || scenario == "parent_death_prctl" {
		wantReason = "signal: terminated"
	}
	require.Equal(t, "shutdown", runEnd.EndReason)
	require.Equal(t, wantReason, runEnd.ShutdownReason)
	require.NotNil(t, runEnd.ShutdownDrained)
	require.True(t, *runEnd.ShutdownDrained)
	require.NotNil(t, runEnd.ExitCode)
	require.EqualValues(t, 255, *runEnd.ExitCode)
	require.True(t, runEnd.Offline)
}
