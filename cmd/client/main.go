package main

import (
	"github.com/xEgorka/project3/cmd/client/internal/app/client"
)

func main() {
	if err := client.Start(); err != nil {
		panic(err)
	}
}
