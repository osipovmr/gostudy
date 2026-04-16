package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}
