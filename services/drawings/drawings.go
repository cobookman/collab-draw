package main

import (
	"encoding/json"
)

type Drawing struct {
	CanvasID string `json:"canvasId"`
}

func (d *Drawing) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Drawing) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, d)
}
