package main

type BenchmarkResult struct {
	getCount int

	missCount int

	setCount int

	StatisticBuckets []BenchmarkStatistic
}
