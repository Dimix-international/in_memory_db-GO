package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMustLoadPath(t *testing.T) {
	t.Parallel()
	cfg := MustLoadPath("../../config/config.yaml")

	assert.Equal(t, "127.0.0.1:3223", cfg.Network.Address)
	assert.Equal(t, time.Minute*5, cfg.Network.IdleTimeout)
	assert.Equal(t, 10, cfg.Network.MaxConnections)
	assert.Equal(t, "4KB", cfg.Network.MaxMessageSize)

	assert.Equal(t, "in_memory", cfg.Engine.Type)

	assert.Equal(t, "local", cfg.Logging.Level)
	assert.Equal(t, "/log/output.log", cfg.Logging.Output)
}

func TestConfigPathEmptyPanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "config path is empty", r)
		} else {
			t.Errorf("The code did not panic")
		}
	}()

	MustLoadConfig()
}

func TestConfigFileNotExistPanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "config file does not exist: invalid/path/config.yaml", r)
		} else {
			t.Errorf("The code did not panic")
		}
	}()

	MustLoadPath("invalid/path/config.yaml")
}
