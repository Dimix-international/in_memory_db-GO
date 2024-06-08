package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
)

func main() {
	log := logger.SetupLogger("")
	reader := bufio.NewReader(os.Stdin)
	//handlerMessage := handler.NewHanlderMessages(log, service.NewParserService(), service.NewAnalyzerService(), db.NewShardMap(ShardValue))

	log.Info("start CLI")

	for {
		request, err := reader.ReadString('\n')
		if err != nil {
			log.Error("error read string", "err", err)
			continue
		}

		request = strings.TrimSpace(request)
		if len(request) == 0 {
			continue
		}

		//handlerMessage.ProcessMessage(request)
	}
}
