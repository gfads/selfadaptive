package mytimer

import (
	"time"
)

type MyTimer struct {
	D            time.Duration
	ChStartTimer chan bool
	ChStopTimer  chan bool
}

func NewMyTimer(d int, chStartTimer, chStopTimer chan bool) MyTimer {
	r := MyTimer{D: time.Duration(d) * time.Second, ChStartTimer: make(chan bool), ChStopTimer: make(chan bool)}
	r.ChStartTimer = chStartTimer
	r.ChStopTimer = chStopTimer
	return r
}

func (t MyTimer) RunMyTimer() {

	for {
		t.ChStartTimer <- true
		time.Sleep(t.D)
		t.ChStopTimer <- true
	}
}
