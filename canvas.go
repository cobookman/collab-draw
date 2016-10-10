package main

import (
	"time"
	"errors"
)

type Canvas struct {
}

type Canvases struct {
	Canvases []Canvas
}


func (c *Canvas) Get(id string) error {
	return errors.New("not completed")
}

func (c *Canvas) Create() error {
	return errors.New("not completed")
}

func (cs *Canvases) GetAll(activeSince time.Time, limit int) error {
	return errors.New("not completed")
}
