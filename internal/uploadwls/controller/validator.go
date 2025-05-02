package controller

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"path/filepath"
	"sort"
	"strings"
)

type ArgsValidator struct {
	args []string
}

func NewArgsValidator(args []string) *ArgsValidator {
	return &ArgsValidator{args: args}
}

func (argsVal *ArgsValidator) Validate() error {
	const expectedArgsLen = 1
	if len(argsVal.args) != expectedArgsLen {
		return fmt.Errorf("accepts %d arg(s), received %d", expectedArgsLen, len(argsVal.args))
	}

	path := argsVal.args[0]
	val := validator.New()
	if errValFile := val.Var(path, "file"); errValFile != nil {
		return fmt.Errorf("file %s does not exist\nError: %w\n", path, errValFile)
	}

	ext := filepath.Ext(path)
	if _, ok := supportedExtensions[ext]; !ok {
		return fmt.Errorf(
			"the file ext %v is not supported. Supported extensions are: %v",
			ext,
			strings.Join(argsVal.getSupportedExtensions(), ", "),
		)
	}

	return nil
}

func (argsVal *ArgsValidator) getSupportedExtensions() []string {
	result := make([]string, 0, len(supportedExtensions))
	for ext, _ := range supportedExtensions {
		result = append(result, ext)
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i]) < strings.ToLower(result[j])
	})

	return result
}
