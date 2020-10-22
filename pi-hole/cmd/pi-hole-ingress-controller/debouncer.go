package main

import (
	"sync"
	"time"
)

type debouncer struct {
	f     func()
	after time.Duration

	mu    sync.Mutex
	timer *time.Timer
}

func (d *debouncer) tick() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, d.f)
}
