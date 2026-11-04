package logging

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// New builds the application logger from LOG_FORMAT ("text" or "json") and LOG_LEVEL ("debug", "info", "warn",
// "error"). Empty values mean text and info.
func New(w io.Writer, format, level string) (*slog.Logger, error) {
	var lvl slog.Level
	if level != "" {
		if err := lvl.UnmarshalText([]byte(level)); err != nil {
			return nil, fmt.Errorf("invalid LOG_LEVEL %q: %w", level, err)
		}
	}
	opts := &slog.HandlerOptions{Level: lvl}

	switch strings.ToLower(format) {
	case "", "text":
		return slog.New(slog.NewTextHandler(w, opts)), nil
	case "json":
		return slog.New(slog.NewJSONHandler(w, opts)), nil
	default:
		return nil, fmt.Errorf("invalid LOG_FORMAT %q: want text or json", format)
	}
}
