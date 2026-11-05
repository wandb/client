package stream

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wandb/wandb/core/internal/observability"
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
	mux.FlushTransactionLogs()
	record, err := r.Read()
	require.NoError(t, err)
	require.EqualValues(t, 42, record.Num)
}
