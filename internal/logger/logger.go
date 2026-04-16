package logger

import (
	"log/slog"
	"os"
)

func Init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}
