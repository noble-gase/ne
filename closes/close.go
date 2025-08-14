package closes

import (
	"fmt"
	"log/slog"
	"sort"
	"sync"
)

type closer struct {
	name     string
	priority int
	fn       func() error
}

var (
	closers []closer
	mutex   sync.Mutex
)

func Add(name string, priority int, fn func() error) {
	mutex.Lock()
	defer mutex.Unlock()

	closers = append(closers, closer{
		name:     name,
		priority: priority,
		fn:       fn,
	})
}

func Close() {
	sort.Slice(closers, func(i, j int) bool {
		return closers[i].priority < closers[j].priority
	})

	for _, closer := range closers {
		fmt.Println("close", closer.name, "...")
		if err := closer.fn(); err != nil {
			slog.Error("close "+closer.name+" failed", "error", err)
		}
	}
}
