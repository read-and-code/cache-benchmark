package cache_client

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HttpCacheClient struct {
	*http.Client

	serverAddress string
}

func (httpCacheClient *HttpCacheClient) get(key string) string {
	response, err := httpCacheClient.Get(httpCacheClient.serverAddress + key)

	if err != nil {
		log.Println(key)

		panic(err)
	}

	if response.StatusCode == http.StatusNotFound {
		return ""
	}

	if response.StatusCode != http.StatusOK {
		panic(response.Status)
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (httpCacheClient *HttpCacheClient) set(key, value string) {
	request, err := http.NewRequest(http.MethodPut, httpCacheClient.serverAddress+key, strings.NewReader(value))

	if err != nil {
		log.Println(err)

		panic(err)
	}

	response, err := httpCacheClient.Do(request)

	if err != nil {
		log.Println(key)

		panic(err)
	}

	if response.StatusCode != http.StatusOK {
		panic(response.Status)
	}
}

func (httpCacheClient *HttpCacheClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		cmd.Value = httpCacheClient.get(cmd.Key)

		return
	}

	if cmd.Name == "set" {
		httpCacheClient.set(cmd.Key, cmd.Value)

		return
	}

	panic("Unknown cmd name " + cmd.Name)
}

func (httpCacheClient *HttpCacheClient) PipelinedRun([]*Cmd) {
	panic("HttpCacheClient pipelined run not implemented")
}

func newHttpCacheClient(host string, port int) *HttpCacheClient {
	httpClient := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}

	return &HttpCacheClient{httpClient, host + ":" + strconv.Itoa(port) + "/cache/"}
}
