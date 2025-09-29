package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v75/github"
	"github.com/turnes/go-github/config"
	api "github.com/turnes/go-github/internal/adapters/http"
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

	// Initialize layers
	githubRepo := repository.NewGithubRepository(client)
	repoService := services.NewRepositoryService(githubRepo)
	httpHandler := api.NewHTTPHandler(repoService)

	r := gin.Default()
	httpHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("API Server starting on port %s...\n", port)
	log.Fatal(r.Run(":" + port))

}
