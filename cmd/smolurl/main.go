package main

import (
	"github.com/mistic0xb/smolurl/internal/config"
)

const PORT = ":8080"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config:" + err.Error())
	}

	cfg.Print()
}
