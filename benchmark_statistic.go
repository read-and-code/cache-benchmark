package main

import "time"

type BenchmarkStatistic struct {
	operationCount int

	duration time.Duration
}
