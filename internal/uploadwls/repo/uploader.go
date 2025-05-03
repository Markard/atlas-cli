package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Markard/atlas-cli/internal/uploadwls/entity"
	"github.com/ctreminiom/go-atlassian/v2/jira/v3"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
	"github.com/ztrue/tracerr"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

type JiraUploader struct {
	url      string
	username string
	token    string
}

func NewJiraUploader(url string, username string, token string) *JiraUploader {
	return &JiraUploader{url: url, username: username, token: token}
}

func (ju JiraUploader) Upload(dailyLogs *entity.DailyWorklogs, ctx context.Context) error {
	jira, err := v3.New(nil, ju.url)
	if err != nil {
		return tracerr.Wrap(fmt.Errorf("check atlassian url\nError: %w\n", err))
	}
	jira.Auth.SetBasicAuth(ju.username, ju.token)

	options := &models.WorklogOptionsScheme{
		Notify:               false,
		AdjustEstimate:       "",
		NewEstimate:          "",
		ReduceBy:             "",
		OverrideEditableFlag: false,
		Expand:               nil,
	}

	errs, ctx := errgroup.WithContext(ctx)

	slog.Info(fmt.Sprintf("[%v] START UPLOAD...", dailyLogs.Date))
	for _, worklog := range dailyLogs.Worklogs {
		if worklog.IssueKey == "" {
			continue
		}

		errs.Go(func() error {
			return ju.upload(worklog, dailyLogs.Date, jira, options)
		})
	}
	uploadErr := errs.Wait()
	slog.Info(fmt.Sprintf("[%v] UPLOAD FINISHED", dailyLogs.Date))

	return uploadErr
}

func (ju JiraUploader) upload(worklog entity.Worklog, date string, jira *v3.Client, options *models.WorklogOptionsScheme) error {
	commentBody := models.CommentNodeScheme{}
	commentBody.Version = 1
	commentBody.Type = "doc"

	payload := &models.WorklogADFPayloadScheme{
		Comment: &models.CommentNodeScheme{
			Version: 1,
			Type:    "doc",
			Content: []*models.CommentNodeScheme{
				{
					Type: "paragraph", Content: []*models.CommentNodeScheme{
						{Type: "text", Text: worklog.Comment},
					},
				},
			},
		},
		Visibility:       nil,
		Started:          worklog.GetStartedAtAsString(),
		TimeSpentSeconds: worklog.GetTimeSpentInSeconds(),
	}

	uploadedWorklog, response, err := jira.Issue.Worklog.Add(
		context.Background(),
		worklog.IssueKey,
		payload,
		options,
	)
	if err != nil && response != nil {
		return fmt.Errorf("Response endpoint: %s, Code:%v\nError: %w\n", response.Endpoint, response.Code,
			tracerr.Wrap(err))
	}

	if response == nil {
		return errors.New("response is nil")
	}

	if uploadedWorklog == nil {
		return errors.New("uploaded worklog is nil")
	}

	slog.Info(ju.getReport(uploadedWorklog, date))
	return nil
}

func (ju JiraUploader) getReport(uw *models.IssueWorklogADFScheme, date string) string {
	var shortComment string
	if len(uw.Comment.Text) > 45 {
		shortComment = uw.Comment.Text[:45] + "..."
	} else {
		shortComment = uw.Comment.Text
	}

	return fmt.Sprintf(
		"[%s] UPLOADED - ID: %s (%s - %-7s - %-7v) (%-29v | %-6v)",
		date,
		uw.ID,
		uw.Started,
		uw.TimeSpent,
		uw.TimeSpentSeconds,
		shortComment,
		uw.IssueID,
	)
}
