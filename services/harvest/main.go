package main

import (
	"time"
)

type TimeEntry struct {
	Project string    `json:"project"`
	Notes   string    `json:"notes"`
	Time    time.Time `json:"time"`
}

func main() {
	// TODO: Get env

	api := NewClient("", "", nil, "")

	// TODO: Pretty print notes for the day
}
