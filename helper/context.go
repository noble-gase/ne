package helper

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
	"google.golang.org/grpc/metadata"
)

func NewCtxWithKeyValue(ctx context.Context, key string, vals ...string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}
	md.Append(key, vals...)
	return metadata.NewIncomingContext(ctx, md)
}

func NewCtxWithTraceId(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}
	if len(md.Get(XTraceId)) != 0 {
		return ctx
	}

	traceId := strings.ReplaceAll(uuid.New().String(), "-", "")

	md.Set(XTraceId, traceId)
	return metadata.NewIncomingContext(ctx, md)
}

func GetValuesFromCtx(ctx context.Context, key string) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md.Get(key)
}

func GetTraceIdFromCtx(ctx context.Context) string {
	vals := GetValuesFromCtx(ctx, XTraceId)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func GetStrValFromCtx(ctx context.Context, key string) string {
	vals := GetValuesFromCtx(ctx, key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func GetBoolValFromCtx(ctx context.Context, key string) bool {
	s := GetStrValFromCtx(ctx, key)

	v, _ := strconv.ParseBool(s)
	return v
}

func GetIntValFromCtx[T constraints.Signed](ctx context.Context, key string) T {
	s := GetStrValFromCtx(ctx, key)

	v, _ := strconv.ParseInt(s, 10, 64)
	return T(v)
}

func GetUintValFromCtx[T constraints.Unsigned](ctx context.Context, key string) T {
	s := GetStrValFromCtx(ctx, key)

	v, _ := strconv.ParseUint(s, 10, 64)
	return T(v)
}

func GetFloatValFromCtx[T constraints.Float](ctx context.Context, key string) T {
	s := GetStrValFromCtx(ctx, key)

	v, _ := strconv.ParseFloat(s, 64)
	return T(v)
}
