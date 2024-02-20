package main

import (
	"time"
)

var DeltaTime = deltaTime{time.Now(), 0}

type deltaTime struct {
	previousTime time.Time
	value        float32
}

func (dt *deltaTime) update() {
	lock.Lock()
	defer lock.Unlock()
	currentTime := time.Now()
	dt.value = float32(currentTime.Sub(dt.previousTime).Seconds())
	dt.previousTime = currentTime
}

func (dt *deltaTime) get() float32 {
	return dt.value
}
