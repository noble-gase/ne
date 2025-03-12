package protos

import (
	"sync"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

// v 默认验证器
var (
	v    protovalidate.Validator
	once sync.Once
)

// Validate 验证消息体
func Validate(msg proto.Message, opts ...protovalidate.ValidatorOption) error {
	if v == nil {
		once.Do(func() {
			v, _ = protovalidate.New(opts...)
		})
	}
	return v.Validate(msg)
}
