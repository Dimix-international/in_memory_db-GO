package network

import (
	"net"
	"testing"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestSercer(t *testing.T) {
	t.Parallel()

	server, err := NewTCPServer(
		&config.NetworkConfig{
			Address:        ":20001",
			MaxConnections: 1,
			IdleTimeout:    time.Second * 5,
			MaxMessageSize: "2KB",
		},
		logger.SetupLogger(""),
	)
	assert.NoError(t, err)

	go func() {
		err = server.Run()
		assert.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	request := "SET weather_2_pm cold_moscow_weather"
	response := "command SET is execute: weather_2_pm"

	connection, err := net.Dial("tcp", "localhost:20001")
	assert.NoError(t, err)

	_, err = connection.Write([]byte(request))
	assert.NoError(t, err)

	buffer := make([]byte, 2048)
	count, err := connection.Read(buffer)
	assert.NoError(t, err)
	assert.Equal(t, []byte(response), buffer[:count])

	//The 2nd connect will wait until the first idleTimeout runs out - check semaphore
	connection2, err := net.Dial("tcp", "localhost:20001")

	request2 := "SET command_2 cold_moscow_weather"
	response2 := "command SET is execute: command_2"

	_, err = connection2.Write([]byte(request2))
	assert.NoError(t, err)

	buffer2 := make([]byte, 2048)
	count, err = connection2.Read(buffer2)

	assert.Equal(t, []byte(response2), buffer2[:count])
}
