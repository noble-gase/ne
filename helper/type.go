package helper

import "fmt"

type M = map[string]any

type X map[string]any

// Size 字节大小
type Size int64

const (
	// B - Byte size
	B Size = 1
	// KiB - KibiByte size
	KiB = 1024 * B
	// MiB - MebiByte size
	MiB = 1024 * KiB
	// GiB - GibiByte size
	GiB = 1024 * MiB
	// TiB - TebiByte size
	TiB = 1024 * GiB
)

// String 实现 Stringer 接口
func (s Size) String() string {
	if s >= TiB {
		return fmt.Sprintf("%.2fTB", float64(s)/float64(TiB))
	}
	if s >= GiB {
		return fmt.Sprintf("%.2fGB", float64(s)/float64(GiB))
	}
	if s >= MiB {
		return fmt.Sprintf("%.2fMB", float64(s)/float64(MiB))
	}
	if s >= KiB {
		return fmt.Sprintf("%.2fKB", float64(s)/float64(KiB))
	}
	return fmt.Sprintf("%dB", s)
}
