package repo

import (
	"errors"
	"fmt"
	"github.com/Markard/atlas-cli/internal/uploadwls/entity"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type YmlProvider struct{}

type rawContent struct {
	Date     string `yaml:"date" validate:"required,datetime=2006-01-02"`
	Worklogs []struct {
		StartedAt string   `yaml:"started_at" validate:"required,datetime=15:04"`
		EndedAt   string   `yaml:"ended_at" validate:"required,datetime=15:04"`
		Comment   string   `yaml:"comment" validate:"required,max=512"`
		IssueKey  string   `yaml:"issue_key" validate:"omitempty,ascii,min=3"`
		Tags      []string `yaml:"tags"`
	} `yaml:"worklogs" validate:"required,dive"`
}

func (p YmlProvider) LoadAndValidate(path string, location *time.Location) (*entity.DailyWorklogs, error) {
	content, errReadFile := os.ReadFile(path)
	if errReadFile != nil {
		return nil, errReadFile
	}

	tmp := new(rawContent)
	errUnmarshal := yaml.Unmarshal([]byte(content), tmp)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	if errVal := p.validateRawContent(tmp); errVal != nil {
		return nil, errVal
	}

	result, errConvert := p.convertToDailyWorklogs(tmp, location)
	if errConvert != nil {
		return nil, errConvert
	}

	return result, nil
}

func (p YmlProvider) validateRawContent(rawContent *rawContent) error {
	if len(rawContent.Worklogs) == 0 {
		return errors.New("no worklogs found")
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if errVal := val.Struct(rawContent); errVal != nil {
		var validationErrors validator.ValidationErrors
		errors.As(errVal, &validationErrors)
		return validationErrors
	}

	return nil
}

func (p YmlProvider) convertToDailyWorklogs(tmp *rawContent, location *time.Location) (*entity.DailyWorklogs, error) {
	result := new(entity.DailyWorklogs)
	result.Date = tmp.Date
	for _, worklog := range tmp.Worklogs {
		tmpWl := new(entity.Worklog)
		tmpWl.IssueKey = worklog.IssueKey
		tmpWl.Comment = worklog.Comment
		tmpWl.Tags = worklog.Tags

		startedAt, errConvertToDatetimeS := p.convertToTime(tmp.Date, worklog.StartedAt, location)
		if errConvertToDatetimeS != nil {
			return nil, errConvertToDatetimeS
		}
		tmpWl.StartedAt = startedAt

		endeAt, errConvertToDatetimeE := p.convertToTime(tmp.Date, worklog.EndedAt, location)
		if errConvertToDatetimeE != nil {
			return nil, errConvertToDatetimeE
		}
		tmpWl.EndedAt = endeAt

		result.Worklogs = append(result.Worklogs, *tmpWl)
	}

	return result, nil
}

func (p YmlProvider) convertToTime(date string, startedAt string, location *time.Location) (*time.Time, error) {
	tmpDatetime := fmt.Sprintf("%s %s", date, startedAt)
	layout := "2006-01-02 15:04"
	datetime, err := time.ParseInLocation(layout, tmpDatetime, location)
	if err != nil {
		return nil, err
	}

	return &datetime, nil
}
