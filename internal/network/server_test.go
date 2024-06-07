package network

import (
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/db"
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/handler"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestSercer(t *testing.T) {
	t.Parallel()

	server, err := NewTCPServer(&config.NetworkConfig{Address: ":20001", MaxConnections: 10, IdleTimeout: time.Second * 10, MaxMessageSize: "2KB"}, &slog.Logger{})
	assert.NoError(t, err)

	log := logger.SetupLogger("")

	handlerRequest := handler.NewHanlderMessages(log, service.NewParserService(), service.NewAnalyzerService(), db.NewShardMap(10))

	server.Run(func(b []byte) {
		handlerRequest.ProcessMessage(string(b))
	})

	request := "hello server"

	connection, err := net.Dial("tcp", "localhost:20001")
	assert.NoError(t, err)

	_, err = connection.Write([]byte(request))
	assert.NoError(t, err)
}
