package main

import (
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

func main() {
	fmt.Println("Hello")
}
