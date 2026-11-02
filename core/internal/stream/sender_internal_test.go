package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/wandb/wandb/core/internal/runhandle"
	"github.com/wandb/wandb/core/internal/runupserter"
	"github.com/wandb/wandb/core/internal/runupsertertest"
	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

func senderWithRun(t *testing.T, entity, project string) *Sender {
	t.Helper()
	runHandle := runhandle.New()
	require.NoError(t, runHandle.Init(runupsertertest.NewTestUpserter(
		t, entity, project, "run-id", runupserter.RunUpserterParams{},
	)))
	return &Sender{runHandle: runHandle}
}

func TestWithRunEntityAndProject(t *testing.T) {
	testCases := []struct {
		name                    string
		entity, project         string
		wantEntity, wantProject string
	}{
		{"both empty", "", "", "run-entity", "run-project"},
		{"empty entity", "", "artifact-project", "run-entity", "artifact-project"},
		{"empty project", "artifact-entity", "", "artifact-entity", "run-project"},
		{"both set", "artifact-entity", "artifact-project", "artifact-entity", "artifact-project"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := senderWithRun(t, "run-entity", "run-project")
			msg := &spb.ArtifactRecord{Entity: tc.entity, Project: tc.project, Name: "table"}
			original := proto.CloneOf(msg)

			got := s.withRunEntityAndProject(msg)

			assert.Equal(t, tc.wantEntity, got.Entity)
			assert.Equal(t, tc.wantProject, got.Project)
			assert.Equal(t, "table", got.Name)
			assert.True(t, proto.Equal(original, msg), "the input record changed")
		})
	}
}

func TestWithRunEntityAndProject_RunNotInitialized(t *testing.T) {
	s := &Sender{runHandle: runhandle.New()}
	msg := &spb.ArtifactRecord{Project: "artifact-project"}

	got := s.withRunEntityAndProject(msg)

	assert.Empty(t, got.Entity)
	assert.Equal(t, "artifact-project", got.Project)
}
