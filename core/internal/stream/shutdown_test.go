package stream

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wandb/wandb/core/internal/observability"
	"github.com/wandb/wandb/core/internal/runwork"
	"github.com/wandb/wandb/core/internal/runworktest"
	"github.com/wandb/wandb/core/internal/settings"
	"github.com/wandb/wandb/core/internal/transactionlog"
	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

func TestShutdownFlushesBufferedTransactionLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "run.wandb")
	w, err := transactionlog.OpenWriter(path)
	require.NoError(t, err)
	defer w.Close()
	require.NoError(t, w.Write(&spb.Record{Num: 42}))
	logger := observability.NewNoOpLogger()
	r, err := transactionlog.OpenReader(path, logger)
	require.NoError(t, err)
	defer r.Close()
	_, err = r.Read()
	require.Error(t, err, "control: the buffered record must not be readable yet")
	require.NoError(t, r.ResetLastRead())
	stream := &Stream{logger: logger}
	stream.transactionWriter.Store((&WriterFactory{Logger: logger}).New(w))
	mux := NewStreamMux()
	require.NoError(t, mux.AddStream("run", stream))
	SyncTransactionLogs(mux.Streams())
	record, err := r.Read()
	require.NoError(t, err)
	require.EqualValues(t, 42, record.Num)
}

// exitWork is Work whose record is a run exit.
type exitWork struct{ runworktest.NoopWork }

func (*exitWork) ToRecord() *spb.Record {
	return &spb.Record{RecordType: &spb.Record_Exit{
		Exit: &spb.RunExitRecord{ExitCode: 3},
	}}
}

func TestWriterSyncsTransactionLogOnExitRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "run.wandb")
	w, err := transactionlog.OpenWriter(path)
	require.NoError(t, err)
	logger := observability.NewNoOpLogger()
	r, err := transactionlog.OpenReader(path, logger)
	require.NoError(t, err)
	defer r.Close()

	// The settings disable the periodic flush, so only the exit record
	// can make these records readable before the writer closes.
	writer := (&WriterFactory{Logger: logger, Settings: settings.New()}).New(w)
	in := make(chan runwork.Work)
	go writer.Do(in)

	in <- runwork.NoRequest(&exitWork{})
	<-writer.Chan()

	record, err := r.Read()
	require.NoError(t, err)
	require.EqualValues(t, 3, record.GetExit().GetExitCode())

	close(in)
	for range writer.Chan() {
	}
}
