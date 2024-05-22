package main

import (
	"fmt"

	"github.com/Dimix-international/in_memory_db-GO/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()

	fmt.Println(cfg)
}
