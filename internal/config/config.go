package config

import (
	"flag"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	once           sync.Once
	ConfigFileName = os.Getenv("CONFIG_FILE_NAME")
)

// Config - a structure containing the project configuration from an env file
type Config struct {
	Engine  *EngineConfig  `yaml:"engine"`
	Network *NetworkConfig `yaml:"network"`
	Logging *LoggingConfig `yaml:"logging"`
	WAL     *WALConfig     `yaml:"wal"`
}

type EngineConfig struct {
	Type string `yaml:"type" env-default:"in_memory"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address" env-default:"127.0.0.1:1111"`
	MaxConnections int           `yaml:"max_connections" env-default:"100"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"1m"`
	MaxMessageSize string        `yaml:"max_message_size" env-default:"2KB"`
}

type LoggingConfig struct {
	Level  string `yaml:"level" env-default:"local"`
	Output string `yaml:"output" env-default:"/log/output.log"`
}

type WALConfig struct {
	FlushingBatchSize    int           `yaml:"flushing_batch_size" env-default:"100"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout" env-default:"10ms"`
	MaxSegmentSize       string        `yaml:"max_segment_size" env-default:"10MB"`
	DataDirectory        string        `yaml:"data_directory" env-default:"/data/venom/wal"`
}

// MustLoadConfig starts reading from the .env file and writing to the configuration structure
func MustLoadConfig() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func fetchConfigPath() string {
	var pathConfig string
	// --config="path/to/config.yaml"
	flag.StringVar(&pathConfig, "config", "", "path to config file")
	flag.Parse()

	if pathConfig == "" {
		godotenv.Load()
		return ConfigFileName
	}

	return pathConfig
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	once.Do(func() {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			panic("cannot read config: " + err.Error())
		}
	})
	return &cfg
}
