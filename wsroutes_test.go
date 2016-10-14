package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"net/url"
	"testing"
)

var (
	wsAddr = flag.String("wsAddr", "localhost:65080", "a host ip address")
)

func TestListenCanvasE2E(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: *wsAddr, Path: "/canvas"}
	_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("dail:", err)
	}
}

func TestHostIp(t *testing.T) {
	i, err := HostIp(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(i.(IpAddr).Ip) == 0 {
		t.Fatal("No ip address :(")
	}
}
