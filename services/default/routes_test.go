package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	httpAddr  = flag.String("httpAddr", "localhost:8080", "a host ip address")
	someDocID string
)

func TestCreateCanvas(t *testing.T) {
	r := httptest.NewRequest("POST", "/api/v1/canvas", nil)
	i, err := CreateCanvas(r)
	if err != nil {
		t.Fatal(err)
	}

	someDocID = i.(Canvas).Id
	if len(i.(Canvas).Id) == 0 {
		t.Fatal("no id for newly created canvas")
	}
}

func TestGetCanvas(t *testing.T) {
	// Should fail w/no id specified
	var r *http.Request
	r = httptest.NewRequest("GET", "/api/v1/canvas", nil)
	if _, err := GetCanvas(r); err == nil {
		t.Fatal("GetCanvas should fail if no id specified in GET Params")
	}

	// Should fail w/a bogus id specified
	r = httptest.NewRequest("GET", "/api/v1/canvas?id=SOMETHING_BOGUS", nil)
	if _, err := GetCanvas(r); err == nil {
		t.Fatal("GetCanvas should fail if bad id specified in GET Params")
	}

	// Should succseed when grabbing existing doc
	if len(someDocID) == 0 {
		TestCreateCanvas(t)
	}
	r = httptest.NewRequest("GET", "/api/v1/canvas?id="+someDocID, nil)
	i, err := GetCanvas(r)
	if err != nil {
		t.Fatal(err)
	}
	if i.(Canvas).Id == someDocID {
		t.Fatalf("Id sent is different that requested (%s, %s)\n", i.(Canvas).Id, someDocID)
	}
}

func TestListCanvas(t *testing.T) {
	// Should Fail with no activeSince specified
	var r *http.Request
	r = httptest.NewRequest("GET", "/api/v1/canvases", nil)
	if _, err := ListCanvases(r); err == nil {
		t.Fatal("Failed to throw error when missing activeSince field")
	}

	// Should Fail with a non RFC3339 activeSince field value
	r = httptest.NewRequest("GET", "/api/v1/canvases?activeSince=1476478142082", nil)
	if _, err := ListCanvases(r); err == nil {
		t.Fatal("Failed to throw error when given non RFC3339 activeSince field")
	}

	// Should Fail with a non integer limit
	r = httptest.NewRequest("GET", "/api/v1/canvases?activeSince="+time.Now().Format(time.RFC3339)+"&limit=F1", nil)
	if _, err := ListCanvases(r); err == nil {
		t.Fatal("Failed to throw error when given a valid RFC3339 & non numeric limit")
	}

	// Should produce a list of canvases when given a valid active since timestamp
	r = httptest.NewRequest("GET", "/api/v1/canvases?activeSince="+time.Unix(0, 0).Format(time.RFC3339), nil)
	i, err := ListCanvases(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(i.(Canvases).Canvases) == 0 {
		t.Fatal("No canvases returned :(")
	}
}

func TestHostIp(t *testing.T) {
	i, err := HostIp(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(i.(IpAddr).Ip) == 0 {
		t.Fatal("no ip address")
	}
}
