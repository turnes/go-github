package cli

import (
	"github.com/turnes/go-github/internal/core/ports"

	"github.com/spf13/cobra"
)

type CLIHandler struct {
	service ports.RepositoryService
}

func NewCLIHandler(service ports.RepositoryService) *CLIHandler {
	return &CLIHandler{service: service}
}

func (c *CLIHandler) RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "git-manager",
		Short: "Git Manager CLI",
	}

	rootCmd.AddCommand(c.RepositoryCommand())

	return rootCmd
}
