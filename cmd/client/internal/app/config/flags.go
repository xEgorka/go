package config

import "flag"

const (
	defaultDBPath  = "gophkeeper.db"
	defaultAppAddr = "https://localhost:8080"
)

var flagDBPath, flagAppAddr string

func parseFlags() {
	flag.StringVar(&flagDBPath, "d", defaultDBPath, "sqlite file path")
	flag.StringVar(&flagAppAddr, "a", defaultAppAddr, "address of application server")
	flag.Parse()
}
