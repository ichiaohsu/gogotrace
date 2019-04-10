package main

import (
	"fmt"
	"time"
)

type reports []Hop

type Hop struct {
	ID      int
	Addr    string
	RTT     time.Duration
	HopTime time.Duration
}

func (h Hop) formatPrint() {
	fmt.Printf("%d %s  %v\n", h.ID, h.Addr, h.RTT)
}

// MaxHopTime print the max time spend between consecutive hops
func (r reports) MaxHopTime() {
	if len(r) == 0 {
		return
	}
	var maxIndex int
	var maxDuration time.Duration
	for i, v := range r {
		if v.HopTime > maxDuration {
			maxDuration = v.HopTime
			maxIndex = i
		}
	}
	fmt.Printf("The longest hop happened at %s . It took %v\n", r[maxIndex].Addr, r[maxIndex].HopTime)
}
