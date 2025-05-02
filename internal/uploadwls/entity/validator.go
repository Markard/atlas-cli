package entity

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type DailyWorklogsValidator struct {
	dailyWorklogs *DailyWorklogs
}

func NewDailyWorklogsValidator(dailyWorklogs *DailyWorklogs) *DailyWorklogsValidator {
	return &DailyWorklogsValidator{dailyWorklogs: dailyWorklogs}
}

func (dailyWorklogsVal *DailyWorklogsValidator) Validate() error {
	if len(dailyWorklogsVal.dailyWorklogs.Worklogs) == 0 {
		return errors.New("no worklogs found")
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if errVal := val.Var(dailyWorklogsVal.dailyWorklogs, "required,dive"); errVal != nil {
		var validationErrors validator.ValidationErrors
		errors.As(errVal, &validationErrors)
		return validationErrors
	}

	return nil
}
