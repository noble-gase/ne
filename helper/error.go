package helper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"

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
	if fn := runtime.FuncForPC(pc); fn != nil {
		attrs = append(attrs, slog.String("caller", fn.Name()))
	}
	attrs = append(attrs, slog.String("location", fmt.Sprintf("%s:%d", file, line)))

	slog.LogAttrs(ctx, slog.LevelError, err.Error(), attrs...)

	return codes.Err
}
