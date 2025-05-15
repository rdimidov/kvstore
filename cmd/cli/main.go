package main

import (
	"context"
	"log"

	"github.com/rdimidov/kvstore/internal/application/config"
	"github.com/rdimidov/kvstore/internal/application/services"
	"github.com/rdimidov/kvstore/internal/infrastructure/storage"
	"github.com/rdimidov/kvstore/internal/presentation/cli"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer config.Cleanup()

	repo := storage.NewMemory()

	compute := services.NewCompute(repo, config.Logger())

	handler, err := services.NewStringHandler(compute)
	if err != nil {
		log.Fatal(err)
	}
	repl := cli.NewCli(handler, nil)

	if err := repl.Run(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
