package closes

import (
	"fmt"
	"log/slog"
	"sort"
	"sync"
)

type closer struct {
	id string
	px Priority
	fn func() error
}

var (
	closers []closer
	mutex   sync.Mutex
)

// Add 将资源关闭操作添加到关闭队列中
//
//	按 [P0 - P100] 顺序关闭（相同优先级按添加顺序关闭）
func Add(id string, px Priority, fn func() error) {
	mutex.Lock()
	defer mutex.Unlock()

	closers = append(closers, closer{
		id: id,
		px: px,
		fn: fn,
	})
}

// Close 关闭队列中的资源
func Close() {
	sort.SliceStable(closers, func(i, j int) bool {
		return closers[i].px < closers[j].px
	})

	for _, v := range closers {
		fmt.Println("close", v.id, "...")
		if err := v.fn(); err != nil {
			slog.Error("close "+v.id+" failed", "error", err)
		}
	}
}
