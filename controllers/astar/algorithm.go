package algorithm

import (
	"rabbitmq/shared"
)

type Astar struct {
}

func (a Astar) Update(vnew float64, vold float64, rold int) (int, float64) {
	rnew := 0
	getnew := 0.0

	if vnew < shared.SV { // The system is in Shut-off Voltage state, task is stopped
		rnew = 0.0
		//fmt.Println("Shut-off voltage state", vnew, vold, rold, rnew)
	} else if vnew < (shared.OV - shared.HYSTERISIS) { // The system is in Low-voltage state, apply AIMD
		if vnew > vold {
			rnew = rold + 1
			//fmt.Println("Low-voltage voltage state (Accelerating)", vnew, vold, rold, rnew)
		} else {
			rnew = rold / 2
			//fmt.Println("Low-voltage voltage state (Reducing)", vnew, vold, rold, rnew)
		}
	} else if vnew > (shared.OV + shared.HYSTERISIS) { // The system is in High Voltage state, apply MIAD
		if vnew < vold {
			rnew = rold - 1
			//fmt.Println("High-voltage state", vnew, vold, rold, rnew)
		} else {
			rnew = rold * 2
			//fmt.Println("High-voltage state", vnew, vold, rold, rnew)
		}
	} else { // The system is at Optimum Voltage state, take no action
		rnew = rold
		//fmt.Println("Optimum voltage state", vnew, vold, rold, rnew)
	}

	// final check of rnew
	if rnew < shared.MINIMUM_TASK_EXECUTION_RATE {
		rnew = shared.MINIMUM_TASK_EXECUTION_RATE
	}
	if rnew > shared.MAXIMUM_TASK_EXECUTION_RATE {
		rnew = shared.MAXIMUM_TASK_EXECUTION_RATE
	}

	if rnew != 0 {
		getnew = 1.0 / float64(rnew)
	} else {
		getnew = 0
	}

	return rnew, getnew
}
