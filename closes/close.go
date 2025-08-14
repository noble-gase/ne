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

	for _, v := range closers {
		fmt.Println("close", v.name, "...")
		if err := v.fn(); err != nil {
			slog.Error("close "+v.name+" failed", "error", err)
		}
	}
}
