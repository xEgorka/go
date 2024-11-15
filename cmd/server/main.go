package main

import "github.com/xEgorka/project3/internal/app/server"

func main() {
	if err := server.Start(); err != nil {
		panic(err)
	}
}
