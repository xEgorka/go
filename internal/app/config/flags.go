package config

import "flag"

const defaultURI = ":8080"

var (
	flagURI   string
	flagDBURI string
)

func parseFlags() {
	flag.StringVar(&flagURI, "a", defaultURI, "server URI")
	flag.StringVar(&flagDBURI, "d", "", "database URI")
	flag.Parse()
}
