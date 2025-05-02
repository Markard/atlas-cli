package entity

import (
	"fmt"
	"time"
)

type Worklog struct {
	StartedAt *time.Time
	EndedAt   *time.Time
	Comment   string
	IssueKey  string
	Tags      []string
}

func (w *Worklog) GetStartedAtAsString() string {
	return w.StartedAt.Format("2006-01-02T15:04:05.000-0700")
}

func (w *Worklog) GetTimeSpentInSeconds() int {
	return int(w.EndedAt.Sub(*w.StartedAt).Seconds())
}

func (w *Worklog) GetReport(date string) string {
	var shortComment string
	if len(w.Comment) > 45 {
		shortComment = w.Comment[:45] + "..."
	} else {
		shortComment = w.Comment
	}

	return fmt.Sprintf(
		"[%v] WL  %s - %s (%-7v | %-29v | %-6v | %v)",
		date,
		w.StartedAt.Format("15:04"),
		w.EndedAt.Format("15:04"),
		w.GetTimeSpentInSeconds(),
		shortComment,
		w.IssueKey,
		w.Tags,
	)
}

type DailyWorklogs struct {
	Date     string
	Worklogs []Worklog
}
