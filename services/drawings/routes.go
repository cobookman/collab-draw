package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func ListenCanvas(r *http.Request, c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func OnIncomingDrawing(drawing Drawing) error {
	return nil
}
