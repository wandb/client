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
			root, err := os.MkdirTemp("", "m2-shutdown-")
			require.NoError(t, err)
			defer os.RemoveAll(root)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
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
			defer log.Close()
			worker.Stdout, worker.Stderr = log, log
			require.NoError(t, worker.Start())
			defer func() { _ = worker.Process.Kill(); _ = worker.Wait() }()
			var ready struct {
				PID  int    `json:"pid"`
				File string `json:"file"`
			}
			require.Eventually(t, func() bool {
				data, err := os.ReadFile(filepath.Join(root, "ready.json"))
				return err == nil && json.Unmarshal(data, &ready) == nil
			}, 20*time.Second, 20*time.Millisecond)
			defer func() { _ = syscall.Kill(ready.PID, syscall.SIGKILL) }()
			if scenario == "parent_death" || scenario == "parent_death_prctl" {
				require.NoError(t, worker.Process.Kill())
			} else {
				sig := syscall.SIGTERM
				if scenario == "sigkill_control" {
					sig = syscall.SIGKILL
				}
				require.NoError(t, syscall.Kill(ready.PID, sig))
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
			reader, err := transactionlog.OpenReader(ready.File, observability.NewNoOpLogger())
			require.NoError(t, err)
			defer reader.Close()
			var rows int
			var exitCode *int32
			for {
				record, err := reader.Read()
				if errors.Is(err, io.EOF) ||
					(scenario == "sigkill_control" && errors.Is(err, io.ErrUnexpectedEOF)) {
					break
				}
				require.NoError(t, err)
				if record.GetHistory() != nil {
					rows++
				}
				if record.GetExit() != nil {
					value := record.GetExit().ExitCode
					exitCode = &value
				}
			}
			if scenario == "sigkill_control" {
				require.Less(t, rows, 25)
				require.Nil(t, exitCode)
			} else {
				require.Equal(t, 25, rows)
				require.NotNil(t, exitCode)
				require.EqualValues(t, 255, *exitCode)
			}
		})
	}
}
