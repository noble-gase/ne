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
	slog.InfoContext(ctx, "[span] time consume",
		slog.String("function", s.n),
		slog.String("duration", time.Since(s.t).String()),
		slog.String("file", s.f),
		slog.Int("line", s.l),
		slog.Any("tags", s.x),
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
		name := fn.Name()
		sp.n = name[strings.Index(name, ".")+1:]
	}
	return sp
}
