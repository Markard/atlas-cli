package controller

import (
	"context"
	"fmt"
	"github.com/Markard/atlas-cli/internal/uploadwls/entity"
	repo2 "github.com/Markard/atlas-cli/internal/uploadwls/repo"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type FileProvider interface {
	LoadAndValidate(path string, location *time.Location) (*entity.DailyWorklogs, error)
}

type Uploader interface {
	Upload(dailyWorklogs *entity.DailyWorklogs, ctx context.Context) error
}

func upload(uploader Uploader, dailyWorklogs *entity.DailyWorklogs, ctx context.Context) error {
	return uploader.Upload(dailyWorklogs, ctx)
}

var supportedExtensions = map[string]struct {
	provider FileProvider
}{
	".yaml": {
		provider: new(repo2.YmlProvider),
	},
	".yml": {
		provider: new(repo2.YmlProvider),
	},
	".json": {},
}

var (
	push bool
	Cmd  = &cobra.Command{
		Use:     "upload-worklogs /path/to/file.yml",
		Aliases: []string{"uw"},
		Short:   "Upload daily worklogs into Jira from JSON or YAML/YML files",
		Long: "Upload daily worklogs into Jira from JSON or YAML/YML files. " +
			"Only logs with the issue_key and the tag work are uploaded.",
		Args: func(cmd *cobra.Command, args []string) error {
			argsVal := NewArgsValidator(args)
			if err := argsVal.Validate(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			errPath := fmt.Errorf("\nPath: %s\n", path)
			fileProvider := supportedExtensions[filepath.Ext(path)].provider
			location, errLocation := time.LoadLocation(os.Getenv("TZ"))
			if errLocation != nil {
				return errLocation
			}

			worklogs, errLoad := fileProvider.LoadAndValidate(path, location)
			if errLoad != nil {
				return fmt.Errorf("%w %w", errLoad, errPath)
			}

			if push {
				uploader := repo2.NewJiraUploader(
					os.Getenv("ATLASSIAN_URL"),
					os.Getenv("ATLASSIAN_USERNAME"),
					os.Getenv("ATLASSIAN_TOKEN"),
				)
				errUpload := upload(uploader, worklogs, context.Background())
				if errUpload != nil {
					return fmt.Errorf("%w %w", errUpload, errPath)
				}
			} else {
				slog.Info(fmt.Sprintf("[%v] READ ONLY", worklogs.Date))
				for _, worklog := range worklogs.Worklogs {
					slog.Info(worklog.GetReport(worklogs.Date))
				}
			}

			return nil
		},
	}
)

func init() {
	Cmd.Flags().BoolVarP(&push, "push", "p", false, "If True, uploads worklogs to Jira. If False, only displays the logs (default: False)")
}
