package test

import "log/slog"

type NullWriter struct{}

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func NewNullLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(NullWriter{}, nil))
}
