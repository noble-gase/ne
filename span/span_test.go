package span

import (
	"context"
	"testing"
)

func TestSpan(t *testing.T) {
	sp := New("hello")
	defer sp.Finish(context.Background())
}
