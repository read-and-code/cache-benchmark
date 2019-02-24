package main

import (
	"fmt"
	"github.com/read-and-code/cache-benchmark/cache_client"
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

func pipeline(cacheClient cache_client.CacheClient, cmds []*cache_client.Cmd, benchmarkResult BenchmarkResult) {
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

func main() {
	fmt.Println("Hello")
}
