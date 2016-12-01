package models

import (
	"golang.org/x/net/context"
	"testing"
)

func TestDrawing_Marshal(t *testing.T) {
	drawing := Drawing{
		ID: "DrawingFoo",
		CanvasID:  "CanvasFooz",
		Points: []Point{
			Point{
				X: 10,
				Y: 5,
			},
		},
	}

	b, err := drawing.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "{\"id\":\"DrawingFoo\",\"canvasId\":\"CanvasFooz\",\"points\":[{\"x\":10,\"y\":5}]}" {
		t.Fatal("drawing did not marshal correctly")
	}
}

func TestDrawing_Unmarshal(t *testing.T) {
	json := "{\"id\":\"DrawingFoo\",\"canvasId\":\"CanvasFooz\",\"points\":[{\"x\":10,\"y\":5}]}"
	drawing := Drawing{}
	if err := drawing.Unmarshal([]byte(json)); err != nil {
		t.Fatal(err)
	}

	if drawing.ID != "DrawingFoo" {
		t.Fatal("DrawingID did not parse correctly")
	}

	if drawing.CanvasID != "CanvasFooz" {
		t.Fatal("CanvasID did not parse correctly")
	}

	if len(drawing.Points) != 1 || drawing.Points[0].X != 10 || drawing.Points[0].Y != 5 {
		t.Fatal("Points did not parse correctly")
	}
}

func TestDrawing_Forward(t *testing.T) {
	ctx := context.Background()

	drawing := Drawing{}
	drawing.AssignID()

	if err := drawing.Forward(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestDrawing_AssignID(t *testing.T) {
	drawing := Drawing{}
	drawing.AssignID()
	if len(drawing.ID) == 0 {
		t.Fatal("Drawing.AssignID failed to generate a unique id")
	}
}
