package controller

import (
	"atlas-cli/internal/getspacekeys/entity"
	"atlas-cli/internal/getspacekeys/repo"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

type Provider interface {
	GetSpaces(batchSize int) ([]*entity.Space, error)
}

func getSpaces(provider Provider, batchSize int) ([]*entity.Space, error) {
	return provider.GetSpaces(batchSize)
}

var (
	batchSize int
	Cmd       = &cobra.Command{
		Use:     "space-keys",
		Aliases: []string{"sk"},
		Short:   "Fetch all spaces in your Confluence instance",
		Long:    "Fetch all spaces in your Confluence instance, displaying their - ids, names and space keys.",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := repo.NewConfluenceProvider(
				os.Getenv("ATLASSIAN_URL"),
				os.Getenv("ATLASSIAN_USERNAME"),
				os.Getenv("ATLASSIAN_TOKEN"),
			)
			spaces, err := getSpaces(provider, batchSize)

			if err != nil {
				return err
			}

			slog.Info("Confluence spaces:")
			for _, space := range spaces {
				slog.Info(space.AsString())
			}

			return nil
		},
	}
)

func init() {
	Cmd.Flags().IntVarP(&batchSize, "batchSize", "b", 100, "Number of spaces to retrieve per page")
}
