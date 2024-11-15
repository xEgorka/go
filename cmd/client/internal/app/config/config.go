// Package config parses flags, gets environment variables and
// provides client configuration parameters.
package config

import (
	"strings"
)

// Config provides client configuration parameters.
type Config struct {
	DBPath        string
	DBDriver      string
	AppAddr       string
	Usr           string
	Token         string
	Key           []byte
	FileSizeLimit int
}

// Setup calculates client configuration parameters.
func Setup() (*Config, error) {
	var cfg Config
	parseFlags()
	cfg.DBPath = flagDBPath
	cfg.AppAddr = flagAppAddr
	if !strings.HasSuffix(cfg.AppAddr, "/") {
		cfg.AppAddr = cfg.AppAddr + "/"
	}
	cfg.DBDriver = "sqlite3"
	cfg.FileSizeLimit = 1 * 1024 * 1024
	return &cfg, nil
}
