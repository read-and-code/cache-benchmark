package main

type BenchmarkConfig struct {
	cacheClientType string

	serverAddress string

	operation string

	totalRequests int

	valueSize int

	totalThreads int

	keyspaceLength int

	pipelineLength int
}
