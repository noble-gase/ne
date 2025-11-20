package sqls

import "strings"

func Minify(sql string) string {
	var (
		out          []rune
		inSingle     bool
		inDouble     bool
		lastWasSpace bool
	)

	for _, r := range sql {
		switch r {
		case '\'':
			// 遇到单引号时切换单引号状态（不在双引号中）
			if !inDouble {
				inSingle = !inSingle
			}
			out = append(out, r)
			lastWasSpace = false
		case '"':
			// 遇到双引号时切换双引号状态（不在单引号中）
			if !inSingle {
				inDouble = !inDouble
			}
			out = append(out, r)
			lastWasSpace = false
		case ' ', '\t', '\n', '\r':
			if inSingle || inDouble {
				// 引号内：保留原样
				out = append(out, r)
			} else {
				// 引号外：压缩为一个空格
				if !lastWasSpace {
					out = append(out, ' ')
					lastWasSpace = true
				}
			}
		default:
			out = append(out, r)
			lastWasSpace = false
		}
	}

	return strings.TrimSpace(string(out))
}
