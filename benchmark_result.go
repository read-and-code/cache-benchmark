package main

import "time"

type BenchmarkResult struct {
	getCount int

	missCount int

	setCount int

	StatisticBuckets []BenchmarkStatistic
}

func (benchmarkResult *BenchmarkResult) addStatistic(bucket int, benchmarkStatistic BenchmarkStatistic) {
	if bucket > len(benchmarkResult.StatisticBuckets)-1 {
		newStatisticBuckets := make([]BenchmarkStatistic, bucket+1)

		copy(newStatisticBuckets, benchmarkResult.StatisticBuckets)

		benchmarkResult.StatisticBuckets = newStatisticBuckets
	}

	statisticBucket := benchmarkResult.StatisticBuckets[bucket]
	statisticBucket.count += benchmarkStatistic.count
	statisticBucket.time += benchmarkStatistic.time

	benchmarkResult.StatisticBuckets[bucket] = statisticBucket
}

func (benchmarkResult *BenchmarkResult) addDuration(duration time.Duration, typeName string) {
	bucket := int(duration / time.Millisecond)

	benchmarkResult.addStatistic(bucket, BenchmarkStatistic{1, duration})

	if typeName == "get" {
		benchmarkResult.getCount++
	} else if typeName == "set" {
		benchmarkResult.setCount++
	} else {
		benchmarkResult.missCount++
	}
}

func (benchmarkResult *BenchmarkResult) addResult(src *BenchmarkResult) {
	for bucket, statistic := range src.StatisticBuckets {
		benchmarkResult.addStatistic(bucket, statistic)
	}

	benchmarkResult.getCount += src.getCount
	benchmarkResult.missCount += src.missCount
	benchmarkResult.setCount += src.setCount
}
