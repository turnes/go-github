package cli

import (
	"context"
	"fmt"

	"github.com/turnes/go-github/internal/core/domain"
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
		Short: "Git Repository Manager CLI",
	}

	rootCmd.AddCommand(c.listCommand())
	rootCmd.AddCommand(c.createCommand())
	rootCmd.AddCommand(c.deleteCommand())
	rootCmd.AddCommand(c.updateCommand())

	return rootCmd
}

func (c *CLIHandler) listCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list [owner]",
		Short: "List all repositories",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repos, err := c.service.ListRepositories(context.Background(), args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Found %d repositories:\n\n", len(repos))
			for _, repo := range repos {
				fmt.Printf("Name: %s\n", repo.Name)
				fmt.Printf("  Description: %s\n", repo.Description)
				fmt.Printf("  Private: %v\n", repo.Private)
				fmt.Printf("  URL: %s\n\n", repo.URL)
			}
			return nil
		},
	}
}

func (c *CLIHandler) createCommand() *cobra.Command {
	var desc string
	var private bool

	cmd := &cobra.Command{
		Use:   "create [owner] [name]",
		Short: "Create a new repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := domain.CreateRepositoryInput{
				Name:        args[1],
				Description: desc,
				Private:     private,
			}

			repo, err := c.service.CreateRepository(context.Background(), args[0], input)
			if err != nil {
				return err
			}

			fmt.Printf("Repository created successfully!\n")
			fmt.Printf("Name: %s\n", repo.Name)
			fmt.Printf("URL: %s\n", repo.URL)
			return nil
		},
	}

	cmd.Flags().StringVarP(&desc, "description", "d", "", "Repository description")
	cmd.Flags().BoolVarP(&private, "private", "p", false, "Make repository private")

	return cmd
}

func (c *CLIHandler) deleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [owner] [name]",
		Short: "Delete a repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.service.DeleteRepository(context.Background(), args[0], args[1]); err != nil {
				return err
			}

			fmt.Printf("Repository %s deleted successfully!\n", args[1])
			return nil
		},
	}
}

func (c *CLIHandler) updateCommand() *cobra.Command {
	var desc string
	var private bool
	var newName string

	cmd := &cobra.Command{
		Use:   "update [owner] [name]",
		Short: "Update a repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := domain.UpdateRepositoryInput{
				Name: args[1],
			}

			if cmd.Flags().Changed("name") {
				input.Name = newName
			}
			if cmd.Flags().Changed("description") {
				input.Description = desc
			}
			if cmd.Flags().Changed("private") {
				input.Private = private
			}

			repo, err := c.service.UpdateRepository(context.Background(), args[0], args[1], input)
			if err != nil {
				return err
			}

			fmt.Printf("Repository updated successfully!\n")
			fmt.Printf("Name: %s\n", repo.Name)
			fmt.Printf("Description: %s\n", repo.Description)
			return nil
		},
	}

	cmd.Flags().StringVarP(&newName, "name", "n", "", "New repository name")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "Repository description")
	cmd.Flags().BoolVarP(&private, "private", "p", false, "Make repository private")

	return cmd
}
