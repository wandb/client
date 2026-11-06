package filestream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgress_AddAcked(t *testing.T) {
	var p Progress
	complete := true

	p.addAcked(&FileStreamRequestJSON{
		Files: map[string]OffsetAndContent{
			"wandb-history.jsonl": {Offset: 0, Content: []string{"a", "b"}},
			"wandb-events.jsonl":  {Offset: 5, Content: nil},
		},
	})
	p.addAcked(&FileStreamRequestJSON{
		Files: map[string]OffsetAndContent{
			"wandb-history.jsonl": {Offset: 2, Content: []string{"c"}},
		},
		Complete: &complete,
	})
	p.addAcked(nil)

	assert.Equal(t, map[string]int{"wandb-history.jsonl": 3}, p.AckedOffsets)
	assert.True(t, p.CompleteAcked)
}

func TestProgress_CloneDoesNotShareOffsets(t *testing.T) {
	p := Progress{AckedOffsets: map[string]int{"x": 1}}
	c := p.clone()
	c.AckedOffsets["x"] = 9
	assert.Equal(t, 1, p.AckedOffsets["x"])
}
