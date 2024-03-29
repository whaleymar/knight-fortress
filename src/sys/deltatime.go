package sys

import (
	"sync"
	"time"
)

var _DELTATIME_LOCK = &sync.Mutex{}
var DeltaTime = deltaTime{time.Now(), 0}

type deltaTime struct {
	previousTime time.Time
	value        float32
}

func (dt *deltaTime) Update() {
	_DELTATIME_LOCK.Lock()
	defer _DELTATIME_LOCK.Unlock()
	currentTime := time.Now()
	dt.value = float32(currentTime.Sub(dt.previousTime).Seconds())
	dt.previousTime = currentTime
}

func (dt *deltaTime) Get() float32 {
	return dt.value
}
