package models

import (
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestNewCanvas(t *testing.T) {
	ctx := context.Background()

	// Test creating canvases
	ncanvas, err := NewCanvas(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Test get of canvas
	ocanvas, err := GetCanvas(ctx, ncanvas.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ocanvas.ID != ncanvas.ID {
		t.Fatal("Failed to grab correct canvas")
	}

	// Test get of canvas when none exist
	if _, err := GetCanvas(ctx, "something-random"); err == nil {
		t.Fatal("Should have thrown ErrNoSuchEntity")
	}

	// Test getting of many canvases
	startTime := time.Now().UTC()
	time.Sleep(100 * time.Millisecond)
	var ncanvases [10]*Canvas
	for i := 0; i < 10; i++ {
		c, err := NewCanvas(ctx)
		if err != nil {
			t.Fatal(err)
		}
		ncanvases[i] = c
	}
	time.Sleep(200 * time.Millisecond)
	canvases, err := GetCanvases(ctx, startTime, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(canvases) != 10 {
		t.Fatalf("Failed to grab canvases, only grabbed %d\n", len(canvases))
	}
}

func TestCanvas_AssignID(t *testing.T) {
	c := Canvas{}
	c.assignID()
	if len(c.ID) == 0 {
		t.Fatal("Failed to assign a unique id to canvas")
	}
}

func TestCanvas_Key(t *testing.T) {
	c := Canvas{}
	c.assignID()
	k := c.Key()
	if k.Name != c.ID {
		t.Fatal("Failed to generate proper key")
	}
}
