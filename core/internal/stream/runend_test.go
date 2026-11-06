package stream

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wandb/wandb/core/internal/filestream"
	"github.com/wandb/wandb/core/internal/observability"
)

func readRunEnd(t *testing.T, dir string) runEndRecord {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, RunEndFileName))
	require.NoError(t, err)
	var rec runEndRecord
	require.NoError(t, json.Unmarshal(data, &rec))
	return rec
}

func newTestRecorder(t *testing.T) (*RunEndRecorder, string) {
	t.Helper()
	dir := t.TempDir()
	return NewRunEndRecorder(dir, observability.NewNoOpLogger()), dir
}

func TestRunEnd_NilRecorderIsNoOp(t *testing.T) {
	var r *RunEndRecorder
	r.Exited(0, false, nil, filestream.Progress{})
	r.UploadStopped(errors.New("x"), filestream.Progress{})
	r.RecordsLost(3)
	r.ShutdownStarted("signal")
	r.ShutdownFinished(true)
	assert.Nil(t, NewRunEndRecorder("", nil))
}

func TestRunEnd_Finished(t *testing.T) {
	r, dir := newTestRecorder(t)

	r.Exited(0, false, nil, filestream.Progress{
		AckedOffsets:  map[string]int{"wandb-history.jsonl": 10},
		CompleteAcked: true,
	})

	rec := readRunEnd(t, dir)
	assert.Equal(t, RunEndFinished, rec.EndReason)
	require.NotNil(t, rec.ExitCode)
	assert.EqualValues(t, 0, *rec.ExitCode)
	assert.True(t, rec.CompleteAcked)
	assert.Equal(t, 10, rec.AckedOffsets["wandb-history.jsonl"])
	assert.Empty(t, rec.UploadError)
}

func TestRunEnd_UploadStoppedIsWrittenBeforeExit(t *testing.T) {
	r, dir := newTestRecorder(t)

	r.UploadStopped(
		errors.New("filestream: failed to upload: 400 Bad Request"),
		filestream.Progress{AckedOffsets: map[string]int{"wandb-history.jsonl": 7}},
	)

	rec := readRunEnd(t, dir)
	assert.Equal(t, RunEndUploadStopped, rec.EndReason)
	assert.Contains(t, rec.UploadError, "400")
	assert.Nil(t, rec.ExitCode, "the run has not exited yet")

	// A later, smaller offset snapshot never moves acked offsets back.
	r.Exited(1, false, nil, filestream.Progress{
		AckedOffsets: map[string]int{"wandb-history.jsonl": 3},
	})
	rec = readRunEnd(t, dir)
	assert.Equal(t, RunEndUploadStopped, rec.EndReason)
	assert.Equal(t, 7, rec.AckedOffsets["wandb-history.jsonl"])
	assert.False(t, rec.CompleteAcked)
}

func TestRunEnd_LostRecordsMarkUploadStopped(t *testing.T) {
	r, dir := newTestRecorder(t)

	r.RecordsLost(5)
	r.Exited(0, false, nil, filestream.Progress{CompleteAcked: true})

	rec := readRunEnd(t, dir)
	assert.Equal(t, RunEndUploadStopped, rec.EndReason)
	assert.EqualValues(t, 5, rec.LostRecords)
}

func TestRunEnd_Shutdown(t *testing.T) {
	r, dir := newTestRecorder(t)

	r.ShutdownStarted("signal: terminated")
	rec := readRunEnd(t, dir)
	assert.Equal(t, RunEndShutdown, rec.EndReason)
	require.NotNil(t, rec.ShutdownDrained)
	assert.False(t, *rec.ShutdownDrained, "not drained until shutdown says so")

	r.Exited(255, false, nil, filestream.Progress{CompleteAcked: true})
	r.ShutdownFinished(true)

	rec = readRunEnd(t, dir)
	assert.Equal(t, RunEndShutdown, rec.EndReason)
	assert.Equal(t, "signal: terminated", rec.ShutdownReason)
	assert.True(t, *rec.ShutdownDrained)
	assert.EqualValues(t, 255, *rec.ExitCode)
}

func TestRunEnd_ShutdownFinishedWithoutStartIsIgnored(t *testing.T) {
	r, dir := newTestRecorder(t)

	r.Exited(0, true, nil, filestream.Progress{})
	r.ShutdownFinished(false)

	rec := readRunEnd(t, dir)
	assert.Equal(t, RunEndFinished, rec.EndReason)
	assert.True(t, rec.Offline)
	assert.Nil(t, rec.ShutdownDrained)
}
