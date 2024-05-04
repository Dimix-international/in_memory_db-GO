package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

var (
	cfg  Config
	once sync.Once
)

type Config struct {
	Env string `env:"ENV" envDefault:"local"`
}

func MustLoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	once.Do(func() {
		if err := env.Parse(&cfg); err != nil {
			fmt.Printf("%+v\n", err)
		}
	})
	return cfg
}
