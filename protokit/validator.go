package protokit

import (
	"sync"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

// v 默认验证器
var (
	v    protovalidate.Validator
	once sync.Once
)

// Validate 验证消息体
func Validate(msg proto.Message) error {
	if v == nil {
		once.Do(func() {
			v, _ = protovalidate.New()
		})
	}
	return v.Validate(msg)
}
