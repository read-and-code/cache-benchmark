package main

import "time"

type BenchmarkResult struct {
	getCount int

	missCount int

	setCount int

	StatisticBuckets []BenchmarkStatistic
}

func (benchmarkResult *BenchmarkResult) addStatistic(index int, benchmarkStatistic BenchmarkStatistic) {
	if index > len(benchmarkResult.StatisticBuckets)-1 {
		newStatisticBuckets := make([]BenchmarkStatistic, index+1)

		copy(newStatisticBuckets, benchmarkResult.StatisticBuckets)

		benchmarkResult.StatisticBuckets = newStatisticBuckets
	}

	statisticBucket := benchmarkResult.StatisticBuckets[index]
	statisticBucket.count += benchmarkStatistic.count
	statisticBucket.duration += benchmarkStatistic.duration

	benchmarkResult.StatisticBuckets[index] = statisticBucket
}

func (benchmarkResult *BenchmarkResult) addDuration(duration time.Duration, typeName string) {
	index := int(duration / time.Millisecond)

	benchmarkResult.addStatistic(index, BenchmarkStatistic{1, duration})

	if typeName == "get" {
		benchmarkResult.getCount++
	} else if typeName == "set" {
		benchmarkResult.setCount++
	} else {
		benchmarkResult.missCount++
	}
}

func (benchmarkResult *BenchmarkResult) addResult(source *BenchmarkResult) {
	for index, benchmarkStatistic := range source.StatisticBuckets {
		benchmarkResult.addStatistic(index, benchmarkStatistic)
	}

	benchmarkResult.getCount += source.getCount
	benchmarkResult.missCount += source.missCount
	benchmarkResult.setCount += source.setCount
}
