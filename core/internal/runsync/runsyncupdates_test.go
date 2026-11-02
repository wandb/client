package runsync_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wandb/wandb/core/internal/runsync"
	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

func artifactRecord(entity, project, runID string) *spb.Record {
	return &spb.Record{RecordType: &spb.Record_Artifact{
		Artifact: &spb.ArtifactRecord{
			Entity:  entity,
			Project: project,
			RunId:   runID,
		},
	}}
}

func TestModify_Artifact_MovesWithTheRun(t *testing.T) {
	updates := &runsync.RunSyncUpdates{
		Entity:  "new-entity",
		Project: "new-project",
		RunID:   "new-id",
	}
	record := artifactRecord("old-entity", "old-project", "old-id")

	updates.Modify(record)

	artifact := record.GetArtifact()
	assert.Equal(t, "new-entity", artifact.Entity)
	assert.Equal(t, "new-project", artifact.Project)
	assert.Equal(t, "new-id", artifact.RunId)
}

func TestModify_Artifact_KeepsWhatIsNotUpdated(t *testing.T) {
	updates := &runsync.RunSyncUpdates{Entity: "new-entity"}
	record := artifactRecord("", "old-project", "old-id")

	updates.Modify(record)

	artifact := record.GetArtifact()
	assert.Equal(t, "new-entity", artifact.Entity)
	assert.Equal(t, "old-project", artifact.Project)
	assert.Equal(t, "old-id", artifact.RunId)
}

func TestModify_NilUpdates_DoesNothing(t *testing.T) {
	var updates *runsync.RunSyncUpdates
	record := artifactRecord("old-entity", "old-project", "old-id")

	updates.Modify(record)

	assert.Equal(t, "old-entity", record.GetArtifact().Entity)
}
