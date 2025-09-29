package main

import (
	"flag"
	"log"

	"github.com/google/go-github/v75/github"
	"github.com/turnes/go-github/config"
	"github.com/turnes/go-github/internal/adapters/cli"
	"github.com/turnes/go-github/internal/adapters/repository"
	"github.com/turnes/go-github/internal/core/services"
)

var (
	env      string
	logLevel string
)

func main() {
	flag.StringVar(&env, "env", ".env", "load .env")
	flag.StringVar(&logLevel, "log", "info", "set log level")
	flag.Parse()

	cfg, err := config.NewFromFile(env)
	if err != nil {
		log.Fatalf("Error loading %s file, %v", env, err)
	}

	client := github.NewClient(nil).WithAuthToken(cfg.Github.PAT)

	githubRepo := repository.NewGithubRepository(client)
	repoService := services.NewRepositoryService(githubRepo)
	cliHandler := cli.NewCLIHandler(repoService)
	cliHandler.RootCommand().Execute()
}
