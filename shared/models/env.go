package models

import (
	"os"
)

// gives us the project hosting this application
func ProjectID() string {
	return os.Getenv("GCLOUD_PROJECT_ID")
}

func UpstreamTopicName() string {
	return os.Getenv("UPSTREAM_DRAWING_TOPIC")
}
