package main

import (
	"fmt"
	"os"

	"github.com/thientran2020/financial-cli/internal/command"
	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	repos, err := repository.NewRepositories(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing repositories: %v\n", err)
		os.Exit(1)
	}

	handler := command.NewHandler(cfg, repos)
	if err := handler.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
