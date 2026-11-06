package stream

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wandb/wandb/core/internal/filestream"
	"github.com/wandb/wandb/core/internal/filestreamtest"
)

// deadFileStream is a FileStream whose uploads stopped with err.
type deadFileStream struct {
	*filestreamtest.FakeFileStream
	dead chan struct{}
	err  error
}

func (f *deadFileStream) Dead() <-chan struct{} { return f.dead }
func (f *deadFileStream) Err() error            { return f.err }
func (f *deadFileStream) Progress() filestream.Progress {
	return filestream.Progress{AckedOffsets: map[string]int{"wandb-history.jsonl": 4}}
}

func TestSenderRecordsUploadStopBeforeExit(t *testing.T) {
	recorder, dir := newTestRecorder(t)
	fake := &deadFileStream{
		FakeFileStream: filestreamtest.NewFakeFileStream(),
		dead:           make(chan struct{}),
		err:            errors.New("filestream: failed to upload: 413 Payload Too Large"),
	}
	s := &Sender{
		fileStream:     fake,
		fileStreamDone: make(chan struct{}),
		runEnd:         recorder,
	}

	close(fake.dead)
	s.recordFileStreamDeath()

	rec := readRunEnd(t, dir)
	assert.Equal(t, RunEndUploadStopped, rec.EndReason)
	assert.Contains(t, rec.UploadError, "413")
	assert.Equal(t, 4, rec.AckedOffsets["wandb-history.jsonl"])
	assert.Nil(t, rec.ExitCode)
}

func TestSenderRecordFileStreamDeathReturnsAfterFinish(t *testing.T) {
	recorder, dir := newTestRecorder(t)
	s := &Sender{
		fileStream:     &deadFileStream{dead: make(chan struct{})},
		fileStreamDone: make(chan struct{}),
		runEnd:         recorder,
	}

	close(s.fileStreamDone)
	s.recordFileStreamDeath()

	assert.NoFileExists(t, dir+"/"+RunEndFileName)
}
