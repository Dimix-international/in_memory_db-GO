package network

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		WAL: &config.WALConfig{
			FlushingBatchSize:    100,
			FlushingBatchTimeout: time.Millisecond * 10,
			MaxSegmentSize:       "10MB",
			DataDirectory:        "/data/venom/wal",
		},
		Network: &config.NetworkConfig{
			Address:        ":20001",
			MaxConnections: 1,
			IdleTimeout:    time.Second * 5,
			MaxMessageSize: "2KB",
		},
	}

	server, err := NewTCPServer(
		cfg,
		logger.SetupLogger(""),
	)
	assert.NoError(t, err)

	go func() {
		err = server.Run(context.TODO())
		assert.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	request := "SET weather_2_pm cold_moscow_weather"
	requestTwo := "SET command_2 2222"
	response := "command SET is execute: weather_2_pm"

	connection, err := net.Dial("tcp", "localhost:20001")
	assert.NoError(t, err)

	_, err = connection.Write([]byte(request))
	assert.NoError(t, err)

	buffer := make([]byte, 2048)
	count, err := connection.Read(buffer)
	assert.NoError(t, err)
	assert.Equal(t, []byte(response), buffer[:count])

	time.Sleep(2 * time.Second)

	_, err = connection.Write([]byte(requestTwo))
	assert.NoError(t, err)

	//The 2nd connect will wait  until the first idleTimeout runs out - check semaphore
	connection2, err := net.Dial("tcp", "localhost:20001")

	request2 := "GET weather_2_pm"
	response2 := "got value from db: cold_moscow_weather"

	request3 := "GET command_2"
	response3 := "got value from db: 2222"

	_, err = connection2.Write([]byte(request2))
	assert.NoError(t, err)

	log.Printf("wait idleTimeout...")

	count, err = connection2.Read(buffer)
	assert.Equal(t, []byte(response2), buffer[:count])

	_, err = connection2.Write([]byte(request3))
	assert.NoError(t, err)

	count, err = connection2.Read(buffer)
	assert.Equal(t, []byte(response3), buffer[:count])

	err = connection2.Close()
	assert.NoError(t, err)
}
