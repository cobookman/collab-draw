package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

type SSEError struct {
	Error   error  `json:"error"`
	Message string `json:"message"`
}

type SSEMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type SSEStatus struct {
	Message string `json:"msg"`
}

type SSEConn struct {
	Writer *http.ResponseWriter
}

type SSEConns struct {
	writers map[string][]*SSEConn
	lock    sync.RWMutex
}

func NewSSEConn(w *http.ResponseWriter) *SSEConn {
	return &SSEConn{
		Writer: w,
	}
}

func NewSSEConns() *SSEConns {
	conns := new(SSEConns)
	conns.writers = make(map[string][]*SSEConn)
	return conns
}

func (c SSEConn) Write(eventType string, data string) error {
	fmt.Fprintf(*c.Writer, "event: %s\ndata: %s\n\n", eventType, data)
	f, ok := (*c.Writer).(http.Flusher)
	if ok {
		f.Flush()
		return nil
	}
	return errors.New("Failed to write sse")
}

func (c SSEConn) WriteInterface(eventType string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.Write(eventType, string(b))
}

func (c SSEConn) WriteErr(err error, msg string) error {
	return c.WriteInterface("error", SSEError{
		Error:   err,
		Message: msg,
	})
}

func (c SSEConn) WriteMsg(v interface{}) error {
	return c.WriteInterface("message", SSEMsg{
		Type: reflect.TypeOf(v).Name(),
		Data: v,
	})
}

func (c SSEConn) WriteAlive() error {
	return c.WriteInterface("health", SSEStatus{
		Message: "Still Alive",
	})
}

func (cs *SSEConns) Add(canvasID string, c *SSEConn) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.writers[canvasID] = append(cs.writers[canvasID], c)
}

func (cs *SSEConns) Remove(canvasID string, c *SSEConn) (int, error) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	canvasWriters := cs.writers[canvasID]
	for i := 0; i < len(canvasWriters); i++ {
		// if we are pointing to same obj, remove the obj from array
		if canvasWriters[i] == c {
			cs.writers[canvasID] = append(canvasWriters[:i], canvasWriters[i+1:]...)
			return len(cs.writers[canvasID]), nil
		}
	}
	return -1, errors.New("Failed to find ResponseWriter in SSEConns")
}

func (cs SSEConns) GetWriters(canvasID string) []*SSEConn {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	writersCopy := make([]*SSEConn, len(cs.writers[canvasID]))
	copy(writersCopy, cs.writers[canvasID])
	return writersCopy
}
