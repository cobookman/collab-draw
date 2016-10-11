package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"testing"
)

var (
	addr = flag.String("addr", "localhost:65080", "http service address")
)

func TestListenCanvasE2E(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/canvas"}
	log.Printf("connecting to %s", u.String())

	_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dail:", err)
	}

}
