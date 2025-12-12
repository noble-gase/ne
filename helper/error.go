package helper

import (
	"context"
	"errors"
	"log/slog"
	"runtime"
	"strings"

	"github.com/noble-gase/ne/codes"
)

type NilError string

func (e NilError) Error() string { return string(e) }

// Error logs the error with caller, then returns the codes.Err
func Error(ctx context.Context, err error, attrs ...slog.Attr) error {
	var code codes.Code
	if errors.As(err, &code) {
		return code
	}

	// Skip level 1 to get the caller function
	pc, file, line, _ := runtime.Caller(1)
	// Get the function details
	var name string
	if fn := runtime.FuncForPC(pc); fn != nil {
		parts := strings.Split(fn.Name(), "/")
		name = parts[len(parts)-1]
	}

	attrs = append(attrs, slog.Group("caller",
		slog.String("func", name),
		slog.String("file", file),
		slog.Int("line", line),
	))

	slog.LogAttrs(ctx, slog.LevelError, err.Error(), attrs...)

	return codes.Err
}
