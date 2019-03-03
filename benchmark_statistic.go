package main

import "time"

type BenchmarkStatistic struct {
	count int

	duration time.Duration
}
