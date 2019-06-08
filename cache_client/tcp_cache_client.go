package cache_client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type TcpCacheClient struct {
	net.Conn
	reader *bufio.Reader
}

func (tcpCacheClient *TcpCacheClient) sendGet(key string) {
	keyLength := len(key)

	_, err := tcpCacheClient.Write([]byte(fmt.Sprintf("G%d %s", keyLength, key)))

	if err != nil {
		log.Println(key)

		panic(err)
	}
}

func (tcpCacheClient *TcpCacheClient) sendSet(key, value string) {
	keyLength := len(key)
	valueLength := len(value)

	_, err := tcpCacheClient.Write([]byte(fmt.Sprintf("S%d %d %s%s", keyLength, valueLength, key, value)))

	if err != nil {
		log.Println(key)

		panic(err)
	}
}

func (tcpCacheClient *TcpCacheClient) sendDelete(key string) {
	keyLength := len(key)

	_, err := tcpCacheClient.Write([]byte(fmt.Sprintf("D%d %s", keyLength, key)))

	if err != nil {
		log.Println(key)

		panic(err)
	}
}

func readLength(reader *bufio.Reader) int {
	text, err := reader.ReadString(' ')

	if err != nil {
		log.Println(text, err)

		return 0
	}

	length, err := strconv.Atoi(strings.TrimSpace(text))

	if err != nil {
		log.Println(text, err)

		return 0
	}

	return length
}

func (tcpCacheClient *TcpCacheClient) receiveResponse() (string, error) {
	valueLength := readLength(tcpCacheClient.reader)

	if valueLength == 0 {
		return "", nil
	}

	if valueLength < 0 {
		errorMessage := make([]byte, -valueLength)
		_, err := io.ReadFull(tcpCacheClient.reader, errorMessage)

		if err != nil {
			return "", err
		}

		return "", errors.New(string(errorMessage))
	}

	value := make([]byte, valueLength)
	_, err := io.ReadFull(tcpCacheClient.reader, value)

	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (tcpCacheClient *TcpCacheClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		tcpCacheClient.sendGet(cmd.Key)

		cmd.Value, cmd.Error = tcpCacheClient.receiveResponse()

		return
	}

	if cmd.Name == "set" {
		tcpCacheClient.sendSet(cmd.Key, cmd.Value)

		_, cmd.Error = tcpCacheClient.receiveResponse()

		return
	}

	if cmd.Name == "delete" {
		tcpCacheClient.sendDelete(cmd.Key)

		_, cmd.Error = tcpCacheClient.receiveResponse()

		return
	}

	panic("Unknown cmd name " + cmd.Name)
}

func (tcpCacheClient *TcpCacheClient) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}

	for _, cmd := range cmds {
		if cmd.Name == "get" {
			tcpCacheClient.sendGet(cmd.Key)
		}

		if cmd.Name == "set" {
			tcpCacheClient.sendSet(cmd.Key, cmd.Value)
		}

		if cmd.Name == "delete" {
			tcpCacheClient.sendDelete(cmd.Key)
		}
	}

	for _, cmd := range cmds {
		cmd.Value, cmd.Error = tcpCacheClient.receiveResponse()
	}
}

func newTcpCacheClient(host string, port int) *TcpCacheClient {
	connection, err := net.Dial("tcp", host+":"+strconv.Itoa(port))

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(connection)

	return &TcpCacheClient{connection, reader}
}
