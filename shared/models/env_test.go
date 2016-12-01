package models

import (
	"testing"
)

func TestEnv(t *testing.T) {
	if len(ProjectID()) == 0 {
		t.Fatal("ProjectID() is empty")
	}

	if len(UpstreamTopicName()) == 0 {
		t.Fatal("UpstreamDrawingTopic() is empty")
	}
}
