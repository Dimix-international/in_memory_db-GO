package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/db"
	"github.com/Dimix-international/in_memory_db-GO/db/storage"
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/handler"
	"github.com/Dimix-international/in_memory_db-GO/internal/models"
	"github.com/Dimix-international/in_memory_db-GO/internal/service"
	"github.com/Dimix-international/in_memory_db-GO/internal/tools"
)

type TCPServer struct {
	maxMessageSize int
	cfg            *config.Config
	semaphore      *tools.Semaphore
	log            *slog.Logger
	storage        *storage.Storage
	idGenerator    *service.IDGenerator
	listener       net.Listener
}

func NewTCPServer(cfg *config.Config, log *slog.Logger) (*TCPServer, error) {
	if log == nil {
		return nil, models.ErrInvalidLogger
	}
	if cfg.Network.MaxConnections <= 0 {
		return nil, models.ErrInvalidMaxConnections
	}

	maxMessageSize, err := tools.ParseSize(cfg.Network.MaxMessageSize)
	if err != nil {
		return nil, err
	}

	wal, err := service.CreateWal(cfg.WAL, log)
	if err != nil {
		return nil, err
	}

	idGenerator := service.NewIDGenerator()

	return &TCPServer{
		maxMessageSize: maxMessageSize,
		cfg:            cfg,
		semaphore:      tools.NewSemaphore(cfg.Network.MaxConnections),
		log:            log,
		storage: storage.NewStorage(
			db.NewDBMap(),
			wal,
			idGenerator,
			log,
		),
		idGenerator: idGenerator,
	}, nil
}

func (s *TCPServer) Run(ctx context.Context) error {
	s.storage.Start()

	listener, err := net.Listen("tcp", s.cfg.Network.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.listener = listener

	var wg sync.WaitGroup
	wg.Add(1)

	go func(ctx context.Context) {
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

			wg.Add(1)
			go func(ctx context.Context, conn net.Conn) {
				defer func() {
					s.semaphore.Release()
					wg.Done()
				}()
				s.semaphore.Acquire()

				s.handleConn(ctx, conn)
			}(ctx, conn)
		}
	}(ctx)

	wg.Wait()

	return nil
}

func (s *TCPServer) Shutdown() error {
	s.storage.Shutdown()
	if err := s.listener.Close(); err != nil {
		return err
	}

	return nil
}

func (s *TCPServer) handleConn(ctx context.Context, conn net.Conn) {
	request := make([]byte, s.maxMessageSize)
	handlerRequest := handler.NewHanlderMessages(service.NewParserService(), service.NewAnalyzerService(), s.storage)

	for {
		if err := conn.SetDeadline(time.Now().Add(s.cfg.Network.IdleTimeout)); err != nil {
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

		ctxWithID := context.WithValue(ctx, models.KeyTxID, s.idGenerator.Generate())
		result := handlerRequest.ProcessMessage(ctxWithID, request[:count])

		if _, err := conn.Write([]byte(result)); err != nil {
			if err != io.EOF {
				s.log.Error("failed to write response", "error", err)
				break
			}
		}
	}

	if err := conn.Close(); err != nil {
		s.log.Error("failed to close connection", "error", err)
	}
}
