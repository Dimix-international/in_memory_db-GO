package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/network"
	"github.com/Dimix-international/in_memory_db-GO/internal/tools"
)

func main() {
	address := flag.String("address", "localhost:3223", "Address of the venom")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	log := logger.SetupLogger("")

	maxMessageSize, err := tools.ParseSize(*maxMessageSizeStr)
	if err != nil {
		log.Error("failed to parse max message size", "err", err)
		os.Exit(1)
	}

	client, err := network.NewTCPClient(*address, maxMessageSize, *idleTimeout)
	if err != nil {
		log.Error("failed to connect with server", "err", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("[venom] >>> ")

		request, err := reader.ReadString('\n')
		if err != nil {
			log.Error("error read string", "err", err)
			continue
		}

		request = strings.TrimSpace(request)
		if len(request) == 0 {
			continue
		}

		resp, err := client.Send([]byte(request))
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				log.Error("connection was closed", "err", err)
			}

			log.Error("failed to send query", "err", err)
		}

		fmt.Println(string(resp))
	}
}
