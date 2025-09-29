package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	Github Github
}

func NewFromFile(file string) (*Config, error) {
	err := godotenv.Load(file)
	if err != nil {
		return nil, err
	}

	return &Config{
		Github: NewGithub(),
	}, nil
}
