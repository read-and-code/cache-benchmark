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

func main() {
	fmt.Println("Hello")
}