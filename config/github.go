package config

import "os"

type Github struct {
	PAT string
}

func NewGithub() Github {
	return Github{
		PAT: os.Getenv("GITHUB_PAT"),
	}
}
