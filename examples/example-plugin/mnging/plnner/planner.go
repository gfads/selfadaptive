/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements a simple planner of MAPE-K.
Date: 28/02/2023
*********************************************************************************/
package plnner

import "main.go/shared"

type Planner struct{}

func NewPlanner() *Planner {
	return &Planner{}
}

func (Planner) Run(fromAnalyser chan shared.ToPlannerChan, toExecutor chan shared.ToPlannerChan) {
	for {
		request := <-fromAnalyser
		toExecutor <- request
	}
}
