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

func pipelinedRun(cacheClient cache_client.CacheClient, cmds []*cache_client.Cmd, benchmarkResult *BenchmarkResult) {
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

func runBenchmark(id int, benchmarkConfig *BenchmarkConfig, channel chan *BenchmarkResult) {
	operationCount := benchmarkConfig.totalRequests / benchmarkConfig.totalThreads
	cacheClient := cache_client.NewCacheClient(benchmarkConfig.cacheClientType, benchmarkConfig.serverAddress)
	cmds := make([]*cache_client.Cmd, 0)
	valuePrefix := strings.Repeat("a", benchmarkConfig.valueSize)
	benchmarkResult := &BenchmarkResult{0, 0, 0, make([]BenchmarkStatistic, 0)}

	for i := 0; i < operationCount; i++ {
		var randomNumber int

		if benchmarkConfig.keyspaceLength > 0 {
			randomNumber = rand.Intn(benchmarkConfig.keyspaceLength)
		} else {
			randomNumber = id*operationCount + i
		}

		key := fmt.Sprintf("%d", randomNumber)
		value := fmt.Sprintf("%s%d", valuePrefix, randomNumber)
		name := benchmarkConfig.operation

		if benchmarkConfig.operation == "mixed" {
			if rand.Intn(2) == 1 {
				name = "set"
			} else {
				name = "get"
			}
		}

		cmd := &cache_client.Cmd{Name: name, Key: key, Value: value, Error: nil}

		if benchmarkConfig.pipelineLength > 1 {
			cmds = append(cmds, cmd)

			if len(cmds) == benchmarkConfig.pipelineLength {
				pipelinedRun(cacheClient, cmds, benchmarkResult)

				cmds = make([]*cache_client.Cmd, 0)
			}
		} else {
			run(cacheClient, cmd, benchmarkResult)
		}
	}

	if len(cmds) != 0 {
		pipelinedRun(cacheClient, cmds, benchmarkResult)
	}

	channel <- benchmarkResult
}

func parseBenchmarkConfig() *BenchmarkConfig {
	benchmarkConfig := &BenchmarkConfig{}

	flag.StringVar(&benchmarkConfig.cacheClientType, "type", "redis", "cache server type")
	flag.StringVar(&benchmarkConfig.serverAddress, "host", "localhost", "cache server address")
	flag.IntVar(&benchmarkConfig.totalRequests, "n", 1000, "total number of requests")
	flag.IntVar(&benchmarkConfig.valueSize, "d", 1000, "data size of SET/GET value in bytes")
	flag.IntVar(&benchmarkConfig.totalThreads, "c", 1, "number of parallel connections")
	flag.StringVar(&benchmarkConfig.operation, "t", "set", "test set, could be get/set/mixed")
	flag.IntVar(&benchmarkConfig.keyspaceLength, "r", 0, "keyspace length, use random keys from 0 to keyspacelen - 1")
	flag.IntVar(&benchmarkConfig.pipelineLength, "p", 1, "pipeline length")
	flag.Parse()

	fmt.Println("cache client type is", benchmarkConfig.cacheClientType)
	fmt.Println("server address is", benchmarkConfig.serverAddress)
	fmt.Println("total number of requests", benchmarkConfig.totalRequests)
	fmt.Println("data size is", benchmarkConfig.valueSize)
	fmt.Println("we have", benchmarkConfig.totalThreads, "parallel connections")
	fmt.Println("operation is", benchmarkConfig.operation)
	fmt.Println("keyspace length is", benchmarkConfig.keyspaceLength)
	fmt.Println("pipeline length is", benchmarkConfig.pipelineLength)

	return benchmarkConfig
}

func main() {
	rand.Seed(time.Now().UnixNano())

	benchmarkConfig := parseBenchmarkConfig()
	channel := make(chan *BenchmarkResult, benchmarkConfig.totalThreads)
	benchmarkResult := &BenchmarkResult{0, 0, 0, make([]BenchmarkStatistic, 0)}
	startTime := time.Now()

	for i := 0; i < benchmarkConfig.totalThreads; i++ {
		go runBenchmark(i, benchmarkConfig, channel)
	}

	for i := 0; i < benchmarkConfig.totalThreads; i++ {
		benchmarkResult.addResult(<-channel)
	}

	duration := time.Now().Sub(startTime)
	operationCount := benchmarkResult.getCount + benchmarkResult.missCount + benchmarkResult.setCount

	fmt.Printf("%d records get\n", benchmarkResult.getCount)
	fmt.Printf("%d records miss\n", benchmarkResult.missCount)
	fmt.Printf("%d records set\n", benchmarkResult.setCount)
	fmt.Printf("%f seconds total\n", duration.Seconds())

	totalOperationCount := 0
	totalDuration := time.Duration(0)

	for i, statisticBucket := range benchmarkResult.StatisticBuckets {
		if statisticBucket.operationCount == 0 {
			continue
		}

		totalOperationCount += statisticBucket.operationCount
		totalDuration += statisticBucket.duration

		fmt.Printf("%d%% requests < %d ms\n", totalOperationCount*100/operationCount, i+1)
	}

	fmt.Printf("%d usec average for eqch request\n", int64(totalDuration/time.Microsecond)/int64(totalOperationCount))
	fmt.Printf("throughput is %f MB/s\n", float64((benchmarkResult.getCount+benchmarkResult.setCount)*benchmarkConfig.valueSize)/1e6/duration.Seconds())
	fmt.Printf("rps is %f\n", float64(operationCount)/float64(duration.Seconds()))
}
