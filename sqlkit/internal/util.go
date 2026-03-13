package internal

import (
	"context"
	"strings"
	"time"
)

type LogFunc = func(ctx context.Context, sql string, cost time.Duration, err error)

var Logger LogFunc

func Minify(sql string) string {
	var (
		inSingle     bool
		inDouble     bool
		lastWasSpace bool
		runes        = []rune(sql)
		n            = len(runes)
	)

	var out strings.Builder
	for i := 0; i < n; i++ {
		r := runes[i]
		switch {
		case !inSingle && !inDouble && r == '-' && i+1 < n && runes[i+1] == '-':
			// 单行注释，跳到行尾
			for i < n && runes[i] != '\n' {
				i++
			}
			// 注释替换为一个空格（避免前后词粘连）
			if !lastWasSpace {
				out.WriteRune(' ')
				lastWasSpace = true
			}

		case !inSingle && !inDouble && r == '/' && i+1 < n && runes[i+1] == '*':
			// 多行注释
			i += 2
			for i+1 < n && (runes[i] != '*' || runes[i+1] != '/') {
				i++
			}
			i++ // 跳过 '/'
			if !lastWasSpace {
				out.WriteRune(' ')
				lastWasSpace = true
			}

		case r == '\'':
			if !inDouble {
				// 处理 '' 转义
				if inSingle && i+1 < n && runes[i+1] == '\'' {
					out.WriteRune(r)
					out.WriteRune(r)
					i++
				} else {
					inSingle = !inSingle
					out.WriteRune(r)
				}
			} else {
				out.WriteRune(r)
			}
			lastWasSpace = false

		case r == '"':
			if !inSingle {
				if inDouble && i+1 < n && runes[i+1] == '"' {
					out.WriteRune(r)
					out.WriteRune(r)
					i++
				} else {
					inDouble = !inDouble
					out.WriteRune(r)
				}
			} else {
				out.WriteRune(r)
			}
			lastWasSpace = false

		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			if inSingle || inDouble {
				out.WriteRune(r)
			} else if !lastWasSpace {
				out.WriteRune(' ')
				lastWasSpace = true
			}

		default:
			out.WriteRune(r)
			lastWasSpace = false
		}
	}
	return strings.TrimSpace(out.String())
}
