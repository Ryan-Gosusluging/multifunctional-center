package config

import "time"

type Settings struct {
    WindowCount    int
    WorkDuration   time.Duration
    MinArrivalTime time.Duration
    MaxArrivalTime time.Duration
    MinServiceTime time.Duration
    MaxServiceTime time.Duration
}

func GetDefault() *Settings {
    return &Settings{
        WindowCount:    5,
        WorkDuration:   15 * time.Second,
        MinArrivalTime: 500 * time.Millisecond,
        MaxArrivalTime: 5 * time.Second,
        MinServiceTime: 200 * time.Millisecond,
        MaxServiceTime: 2 * time.Second,
    }
}