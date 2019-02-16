package cache_client

type CacheClient interface {
	Run(*Cmd)

	PipelinedRun([]*Cmd)
}

func NewCacheClient(cacheClientType, serverAddress string) CacheClient {
	if cacheClientType == "http" {
		return newHttpCacheClient(serverAddress)
	}

	panic("Unknown cache client type " + cacheClientType)
}
