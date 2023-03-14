package onoff

import (
	"selfadaptive/shared"
	"time"
)

type OnOff struct{}

func (OnOff) Update(d time.Duration, g time.Duration) time.Duration {
	if d.Milliseconds() > g.Milliseconds() {
		return time.Duration(shared.MinOnoff * time.Millisecond)
	} else {
		return time.Duration(shared.MaxOnoff * time.Millisecond)
	}
}
