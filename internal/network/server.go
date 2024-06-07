package network

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/models"
	"github.com/Dimix-international/in_memory_db-GO/internal/tools"
)

type TCPServer struct {
	maxMessageSize int
	cfg            *config.NetworkConfig
	semaphore      *tools.Semaphore
	log            *slog.Logger
}

// type TCPHandler = func([]byte) []byte
type TCPHandler = func([]byte)

func NewTCPServer(cfg *config.NetworkConfig, log *slog.Logger) (*TCPServer, error) {
	if log == nil {
		return nil, models.ErrInvalidLogger
	}
	if cfg.MaxConnections <= 0 {
		return nil, models.ErrInvalidMaxConnections
	}

	maxMessageSize, err := tools.ParseSize(cfg.MaxMessageSize)
	if err != nil {
		return nil, err
	}

	return &TCPServer{
		maxMessageSize: maxMessageSize,
		cfg:            cfg, semaphore: tools.NewSemaphore(cfg.MaxConnections),
		log: log,
	}, nil
}

func (s *TCPServer) Run(handler TCPHandler) error {
	listener, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.log.Error("failed to accept", "error", err)
				continue
			}

			go func(conn net.Conn) {
				defer s.semaphore.Release()
				s.semaphore.Acquire()

				s.handleConn(conn, handler)
			}(conn)
		}
	}()

	wg.Wait()

	if err := listener.Close(); err != nil {
		s.log.Error("failed to close listener", "error", err)
	}

	return nil
}

func (s *TCPServer) handleConn(conn net.Conn, handler TCPHandler) {
	request := make([]byte, s.maxMessageSize)

	for {
		if err := conn.SetDeadline(time.Now().Add(s.cfg.IdleTimeout)); err != nil {
			s.log.Error("failed to set read deadline", "error", err)
			break
		}

		count, err := conn.Read(request)
		if err != nil {
			if err != io.EOF {
				s.log.Error("failed to read", "error", err)
			}

			break
		}

		fmt.Println(request)
		handler(request[:count])
	}

	if err := conn.Close(); err != nil {
		s.log.Error("failed to close connection", "error", err)
	}
}
