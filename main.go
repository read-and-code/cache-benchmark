package main

import (
	"flag"
	"fmt"
	"github.com/read-and-code/cache-benchmark/cache_client"
	"math/rand"
	"strings"
	"time"
)

func run(cacheClient cache_client.CacheClient, cmd *cache_client.Cmd, benchmarkResult *BenchmarkResult) {
	expectedValue := cmd.Value
	startTime := time.Now()

	cacheClient.Run(cmd)

	duration := time.Now().Sub(startTime)
	resultType := cmd.Name

	if resultType == "get" {
		if cmd.Value == "" {
			resultType = "miss"
		} else if cmd.Value != expectedValue {
			panic(cmd)
		}
	}

	benchmarkResult.addDuration(duration, resultType)
}

func pipeline(cacheClient cache_client.CacheClient, cmds []*cache_client.Cmd, benchmarkResult *BenchmarkResult) {
	expectedValues := make([]string, len(cmds))

	for i, cmd := range cmds {
		if cmd.Name == "get" {
			expectedValues[i] = cmd.Value
		}
	}

	startTime := time.Now()

	cacheClient.PipelinedRun(cmds)

	duration := time.Now().Sub(startTime)

	for i, cmd := range cmds {
		resultType := cmd.Name

		if resultType == "get" {
			if cmd.Value == "" {
				resultType = "miss"
			} else if cmd.Value != expectedValues[i] {
				fmt.Println(expectedValues[i])

				panic(cmd.Value)
			}
		}

		benchmarkResult.addDuration(duration, resultType)
	}
}

func operate(id, count int, channel chan *BenchmarkResult) {
	cacheClient := cache_client.NewCacheClient(cacheClientType, serverAddress)
	cmds := make([]*cache_client.Cmd, 0)
	valuePrefix := strings.Repeat("a", valueSize)
	benchmarkResult := &BenchmarkResult{0, 0, 0, make([]BenchmarkStatistic, 0)}

	for i := 0; i < count; i++ {
		var randomNumber int

		if keyspaceLength > 0 {
			randomNumber = rand.Intn(keyspaceLength)
		} else {
			randomNumber = id*count + i
		}

		key := fmt.Sprintf("%d", randomNumber)
		value := fmt.Sprintf("%s%d", valuePrefix, randomNumber)
		name := operation

		if operation == "mixed" {
			if rand.Intn(2) == 1 {
				name = "set"
			} else {
				name = "get"
			}
		}

		cmd := &cache_client.Cmd{Name: name, Key: key, Value: value, Error: nil}

		if pipelineLength > 1 {
			cmds = append(cmds, cmd)

			if len(cmds) == pipelineLength {
				pipeline(cacheClient, cmds, benchmarkResult)

				cmds = make([]*cache_client.Cmd, 0)
			}
		} else {
			run(cacheClient, cmd, benchmarkResult)
		}
	}

	if len(cmds) != 0 {
		pipeline(cacheClient, cmds, benchmarkResult)
	}

	channel <- benchmarkResult
}

var cacheClientType, serverAddress, operation string
var totalRequests, valueSize, threads, keyspaceLength, pipelineLength int

func init() {
	flag.StringVar(&cacheClientType, "type", "redis", "cache server type")
	flag.StringVar(&serverAddress, "host", "localhost", "cache server address")
	flag.IntVar(&totalRequests, "n", 1000, "total number of requests")
	flag.IntVar(&valueSize, "d", 1000, "data size of SET/GET value in bytes")
	flag.IntVar(&threads, "c", 1, "number of parallel connections")
	flag.StringVar(&operation, "t", "set", "test set, could be get/set/mixed")
	flag.IntVar(&keyspaceLength, "r", 0, "keyspace length, use random keys from 0 to keyspacelen - 1")
	flag.IntVar(&pipelineLength, "p", 1, "pipeline length")
	flag.Parse()

	fmt.Println("cache client type is", cacheClientType)
	fmt.Println("server address is", serverAddress)
	fmt.Println("total number of requests", totalRequests)
	fmt.Println("data size is", valueSize)
	fmt.Println("we have", threads, "parallel connections")
	fmt.Println("operation is", operation)
	fmt.Println("keyspace length is", keyspaceLength)
	fmt.Println("pipeline length is", pipelineLength)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	channel := make(chan *BenchmarkResult, threads)
	benchmarkResult := &BenchmarkResult{0, 0, 0, make([]BenchmarkStatistic, 0)}
	startTime := time.Now()

	for i := 0; i < threads; i++ {
		go operate(i, totalRequests/threads, channel)
	}

	for i := 0; i < threads; i++ {
		benchmarkResult.addResult(<-channel)
	}

	duration := time.Now().Sub(startTime)
	totalCount := benchmarkResult.getCount + benchmarkResult.missCount + benchmarkResult.setCount

	fmt.Printf("%d records get\n", benchmarkResult.getCount)
	fmt.Printf("%d records miss\n", benchmarkResult.missCount)
	fmt.Printf("%d records set\n", benchmarkResult.setCount)
	fmt.Printf("%f seconds total\n", duration.Seconds())

	statisticCountSum := 0
	statisticTimeSum := time.Duration(0)

	for i, statisticBucket := range benchmarkResult.StatisticBuckets {
		if statisticBucket.count == 0 {
			continue
		}

		statisticCountSum += statisticBucket.count
		statisticTimeSum += statisticBucket.time

		fmt.Printf("%d%% requests < %d ms\n", statisticCountSum*100/totalCount, i+1)
	}

	fmt.Printf("%d usec average for eqch request\n", int64(statisticTimeSum/time.Microsecond)/int64(statisticCountSum))
	fmt.Printf("throughput is %f MB/s\n", float64((benchmarkResult.getCount+benchmarkResult.setCount)*valueSize)/1e6/duration.Seconds())
	fmt.Printf("rps is %f\n", float64(totalCount)/float64(duration.Seconds()))
}
