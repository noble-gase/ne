package helper

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/noble-gase/ne/codes"
)

type demo struct {
	ID   int
	Name string
}

func TestError(t *testing.T) {
	ctx := context.Background()
	err := errors.New("oh no")
	_ = Error(ctx, codes.FromError(err))
	_ = Error(ctx, err)
	_ = Error(ctx, err, slog.Int("id", 1), slog.String("name", "hello"))
	_ = Error(ctx, err, slog.Any("fn", func() {}))
	_ = Error(ctx, err, slog.Any("demo", &demo{ID: 1, Name: "hello"}))
}
