package main

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
