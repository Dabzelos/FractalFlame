package io

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/errors"
)

type Output struct {
	w      io.Writer
	logger slog.Logger
}

func NewWriter(w io.Writer, logger *slog.Logger) *Output {
	return &Output{w: w, logger: *logger}
}

func (o *Output) Write(messages ...interface{}) {
	message := fmt.Sprintln(messages...)

	_, err := o.w.Write([]byte(message))
	if err != nil {
		o.logger.Error("output error occurred", errors.ErrOutPut{}.Error(), err)
	}
}

type Input struct {
	r io.Reader
}
