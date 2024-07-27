package resteasy

import (
	"log/slog"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
