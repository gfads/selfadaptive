package onoff

import "time"

type OnOff struct{}

func (OnOff) Update(d time.Duration, g time.Duration) time.Duration {
	if d.Milliseconds() > g.Milliseconds() {
		//fmt.Println("Acima da Meta")
		return time.Duration(10 * time.Millisecond)
	} else {
		//fmt.Println("Abaixo da Meta")
		return time.Duration(600 * time.Millisecond)
	}
}
