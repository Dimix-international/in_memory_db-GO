package network

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTCPCli(t *testing.T) {
	t.Parallel()

	request := "SET weather_2_pm cold_moscow_weather"
	response := "command SET is execute: weather_2_pm"

	listener, err := net.Listen("tcp", ":10001")
	assert.NoError(t, err)

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}

		buffer := make([]byte, 2048)
		count, err := connection.Read(buffer)
		assert.NoError(t, err)
		assert.Equal(t, buffer[:count], []byte(request))

		_, err = connection.Write([]byte(response))
		assert.NoError(t, err)

		defer func() {
			err = connection.Close()
			assert.NoError(t, err)
		}()
	}()

	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		client, err := NewTCPClient("127.0.0.1:10001", 2048, time.Minute)
		assert.NoError(t, err)

		resp, err := client.Send([]byte(request))
		assert.NoError(t, err)
		assert.Equal(t, resp, []byte(response))
		wg.Done()
	}()

	go func() {
		client, err := NewTCPClient("127.0.0.1:10001", 2048, time.Millisecond*50)
		assert.NoError(t, err)

		_, err = client.Send([]byte(request))
		assert.Error(t, err)
		wg.Done()
	}()

	wg.Wait()
}
