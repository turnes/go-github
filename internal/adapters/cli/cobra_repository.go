package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/turnes/go-github/internal/core/domain"

	"github.com/spf13/cobra"
)

func (c *CLIHandler) RepositoryCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "repo",
		Short: "Repository commands",
	}

	rootCmd.AddCommand(c.listCommand())
	rootCmd.AddCommand(c.getCommand())
	rootCmd.AddCommand(c.createCommand())
	rootCmd.AddCommand(c.deleteCommand())
	rootCmd.AddCommand(c.updateCommand())
	rootCmd.AddCommand(c.setPrivateCommand())

	return rootCmd
}

func (c *CLIHandler) getCommand() *cobra.Command {
	var owner string
	var name string

	return &cobra.Command{
		Use:   "get [owner] [name]",
		Short: "get a repositories",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner = args[0]
			name = args[1]
			repo, err := c.service.GetRepository(context.Background(), owner, name)
			if err != nil {
				return err
			}

			fmt.Printf("Name: %s\n", repo.Name)
			fmt.Printf("  Description: %s\n", repo.Description)
			fmt.Printf("  Private: %v\n", repo.Private)
			fmt.Printf("  URL: %s\n\n", repo.URL)

			return nil
		},
	}
}

func (c *CLIHandler) listCommand() *cobra.Command {
	var owner string
	var sort string
	var direction string
	cmd := &cobra.Command{
		Use:   "list [owner]",
		Short: "List all repositories",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner = args[0]

			repos, err := c.service.ListRepositories(context.Background(), owner, sort, direction)
			if err != nil {
				return err
			}

			fmt.Printf("Found %d repositories:\n\n", len(repos))
			for _, repo := range repos {
				fmt.Printf("Name: %s\n", repo.Name)
				fmt.Printf("  Description: %s\n", repo.Description)
				fmt.Printf("  Private: %v\n", repo.Private)
				fmt.Printf("  URL: %s\n", repo.URL)
				fmt.Printf("  Updated: %s\n", repo.UpdatedAt)
				fmt.Printf("  Pushed: %s\n\n", repo.PushedAt)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&sort, "sort", "s", "created", "sort the results by created, updated, pushed, full_name")
	cmd.Flags().StringVarP(&direction, "direction", "d", "desc", "order to sort by asc, desc")

	return cmd
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

func (c *CLIHandler) setPrivateCommand() *cobra.Command {
	var owner string
	var name string
	var private bool
	var public bool
	var visibility bool

	cmd := &cobra.Command{
		Use:   "set-visibility [owner] [name]",
		Short: "Set a repository private or public",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner = args[0]
			name = args[1]

			if cmd.Flags().Changed("private") && cmd.Flags().Changed("public") {
				return errors.New("only one of private or public can be set")
			}
			if cmd.Flags().Changed("private") {
				visibility = true
			}
			if cmd.Flags().Changed("public") {
				visibility = false
			}

			repo, err := c.service.SetPrivate(context.Background(), owner, name, visibility)
			if err != nil {
				return err
			}

			fmt.Printf("Repository visibility updated successfully!\n")
			fmt.Printf("Name: %s\n", repo.Name)
			fmt.Printf("Private: %t\n", repo.Private)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&private, "private", "", false, "Make repository private")
	cmd.Flags().BoolVarP(&public, "public", "", false, "Make repository public")
	cmd.MarkFlagsOneRequired("private", "public")

	return cmd
}
