package cache_client

type CacheClient interface {
	Run(*Cmd)

	PipelinedRun([]*Cmd)
}

func NewCacheClient(cacheClientType, host string, port int) CacheClient {
	if cacheClientType == "http" {
		return newHttpCacheClient(host, port)
	}

	panic("Unknown cache client type " + cacheClientType)
}
