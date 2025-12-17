package span

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
	"time"
)

// Span 记录耗时
type Span struct {
	f string
	l int
	n string
	t time.Time
	x []string
}

// Finish 记录耗时
func (s *Span) Finish(ctx context.Context) {
	slog.LogAttrs(ctx, slog.LevelInfo, "[span] time consume",
		slog.String("duration", time.Since(s.t).String()),
		slog.Any("tags", s.x),
		slog.Attr{
			Key: "caller",
			Value: slog.GroupValue(
				slog.String("func", s.n),
				slog.String("file", s.f),
				slog.Int("line", s.l),
			),
		},
	)
}

// New 返回一个 span 记录耗时
//
// Example:
//
//	sp := span.New()
//	defer sp.Finish(ctx)
func New(tags ...string) *Span {
	sp := &Span{
		t: time.Now(),
		x: tags,
	}
	// Skip level 1 to get the caller function
	pc, file, line, _ := runtime.Caller(1)
	sp.f, sp.l = file, line
	// Get the function details
	if fn := runtime.FuncForPC(pc); fn != nil {
		parts := strings.Split(fn.Name(), "/")
		sp.n = parts[len(parts)-1]
	}
	return sp
}
