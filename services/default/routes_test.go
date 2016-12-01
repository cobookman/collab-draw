package main

import (
	"github.com/cobookman/collabdraw/shared/models"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateCanvas(t *testing.T) {
	r := httptest.NewRequest("POST", "/api/v1/canvas", nil)
	i, err := HandleCreateCanvas(r)
	if err != nil {
		t.Fatal(err)
	}

	if len(i.(*models.Canvas).ID) == 0 {
		t.Fatal("no id for newly created canvas")
	}
}

func TestGetCanvas(t *testing.T) {
	// Should fail w/no id specified
	var r *http.Request
	r = httptest.NewRequest("GET", "/api/v1/canvas", nil)
	if _, err := HandleGetCanvas(r); err == nil {
		t.Fatal("GetCanvas should fail if no id specified in GET Params")
	}

	// Should fail w/a bogus id specified
	r = httptest.NewRequest("GET", "/api/v1/canvas?id=SOMETHING_BOGUS", nil)
	if _, err := HandleGetCanvas(r); err == nil {
		t.Fatal("GetCanvas should fail if bad id specified in GET Params")
	}

	// Should succseed when grabbing existing doc
	canvas, err := models.NewCanvas(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	r = httptest.NewRequest("GET", "/api/v1/canvas?id="+canvas.ID, nil)
	i, err := HandleGetCanvas(r)
	if err != nil {
		t.Fatal(err)
	}
	if i.(*models.Canvas).ID != canvas.ID {
		t.Fatalf("Id sent is different that requested (%s, %s)\n", i.(*models.Canvas).ID, canvas.ID)
	}
}

func TestListCanvas(t *testing.T) {
	// Should Fail with no activeSince specified
	var r *http.Request
	r = httptest.NewRequest("GET", "/api/v1/canvases", nil)
	if _, err := HandleListCanvases(r); err == nil {
		t.Fatal("Failed to throw error when missing createdSince field")
	}

	// Should Fail with a non RFC3339 activeSince field value
	r = httptest.NewRequest("GET", "/api/v1/canvases?createdSince=1476478142082", nil)
	if _, err := HandleListCanvases(r); err == nil {
		t.Fatal("Failed to throw error when given non RFC3339 activeSince field")
	}

	// Should Fail with a non integer limit
	r = httptest.NewRequest("GET", "/api/v1/canvases?createdSinace="+time.Now().Format(time.RFC3339)+"&limit=F1", nil)
	if _, err := HandleListCanvases(r); err == nil {
		t.Fatal("Failed to throw error when given a valid RFC3339 & non numeric limit")
	}

	// Should produce a list of canvases when given a valid active since timestamp
	for i := 0; i < 10; i++ {
		_, err := models.NewCanvas(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	}
	r = httptest.NewRequest("GET", "/api/v1/canvases?createdSince="+time.Unix(0, 0).Format(time.RFC3339), nil)
	i, err := HandleListCanvases(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(i.([]*models.Canvas)) == 0 {
		t.Fatal("No canvases returned")
	}
}
