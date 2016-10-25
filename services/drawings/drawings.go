package main

import (
	"encoding/json"
)

var (
	dc          *datastore.Client
	DrawingKind = "Drawing"
)

func init() {
	ctx := context.Background()
	projectID := os.Getenv("GCLOUD_DATASET_ID")

	var err error
	dc, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
}

type Drawing struct {
	CanvasID  string `json:"canvasId"`
	DrawingID string `json:"drawingId"`
}

func (d *Drawing) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Drawing) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, d)
}

func GetDrawings(ctx context.Context, canvas Canvas) (*[]Drawing, error) {
	q := datastore.NewQuery(DrawingKind).
		Filter("CanvasID =", canvas.ID)
	drawings := &make([]Drawing)
	keys, err := dc.GetAll(ctx, q, drawings)
	if err != nil {
		return nil, err
	}

	// Attach Drawing Key Encoding to each obj for future ref.
	for i, key := range keys {
		drawings[i].DrawingID = key.Encode()
	}
	return drawings, nil
}
