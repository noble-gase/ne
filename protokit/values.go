package protokit

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dromara/carbon/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ErrUnsupportedMessage 当尝试把扁平字符串解析成非 WKT 的 message 类型时返回。
// 调用方可以用 errors.Is 来区分"类型不支持"和"格式错误"。
var ErrUnsupportedMessage = errors.New("unsupported message type")

// MessageToValues parses proto.Message into url.Values
func MessageToValues(msg proto.Message) url.Values {
	return parseMessage(msg.ProtoReflect())
}

func parseMessage(msg protoreflect.Message) url.Values {
	query := url.Values{}
	fields := msg.Descriptor().Fields()

	// 用 fields.Get(i) 而不是 ByNumber，避免 field number 不连续导致的性能问题
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		key := string(fd.Name())

		// 跳过未设置的字段，避免默认值（0, false, ""）污染 query string
		if !fd.IsList() && !fd.IsMap() && !msg.Has(fd) {
			continue
		}

		val := msg.Get(fd)

		// list
		if fd.IsList() {
			for k, v := range listToValues(fd, val.List()) {
				query[k] = append(query[k], v...)
			}
			continue
		}

		// map
		if fd.IsMap() {
			for k, v := range mapToValues(fd, val.Map()) {
				newKey := key + "." + k
				query[newKey] = append(query[newKey], v...)
			}
			continue
		}

		// 其他类型
		switch fd.Kind() {
		case protoreflect.StringKind:
			query.Add(key, val.String())
		case protoreflect.BoolKind:
			query.Add(key, strconv.FormatBool(val.Bool()))
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			query.Add(key, strconv.FormatInt(val.Int(), 10))
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			query.Add(key, strconv.FormatUint(val.Uint(), 10))
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			query.Add(key, strconv.FormatFloat(val.Float(), 'f', -1, 64))
		case protoreflect.BytesKind:
			query.Add(key, base64.URLEncoding.EncodeToString(val.Bytes()))
		case protoreflect.EnumKind:
			query.Add(key, enumToString(fd.Enum(), val.Enum()))
		case protoreflect.MessageKind, protoreflect.GroupKind:
			// 先尝试 WKT 扁平化，保持与 getProtoMessage 对称
			if s, ok := wktToString(val.Message()); ok {
				query.Add(key, s)
				break
			}
			for k, v := range parseMessage(val.Message()) {
				newKey := key + "." + k
				query[newKey] = append(query[newKey], v...)
			}
		}
	}
	return query
}

func listToValues(fd protoreflect.FieldDescriptor, list protoreflect.List) url.Values {
	query := url.Values{}
	if list.Len() == 0 {
		return query
	}

	key := string(fd.Name())
	switch fd.Kind() {
	case protoreflect.StringKind:
		for i := 0; i < list.Len(); i++ {
			query.Add(key, list.Get(i).String())
		}
	case protoreflect.BoolKind:
		for i := 0; i < list.Len(); i++ {
			query.Add(key, strconv.FormatBool(list.Get(i).Bool()))
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		for i := 0; i < list.Len(); i++ {
			query.Add(key, strconv.FormatInt(list.Get(i).Int(), 10))
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		for i := 0; i < list.Len(); i++ {
			query.Add(key, strconv.FormatUint(list.Get(i).Uint(), 10))
		}
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		for i := 0; i < list.Len(); i++ {
			query.Add(key, strconv.FormatFloat(list.Get(i).Float(), 'f', -1, 64))
		}
	case protoreflect.BytesKind:
		for i := 0; i < list.Len(); i++ {
			query.Add(key, base64.URLEncoding.EncodeToString(list.Get(i).Bytes()))
		}
	case protoreflect.EnumKind:
		for i := 0; i < list.Len(); i++ {
			query.Add(key, enumToString(fd.Enum(), list.Get(i).Enum()))
		}
	case protoreflect.MessageKind, protoreflect.GroupKind:
		for i := 0; i < list.Len(); i++ {
			msg := list.Get(i).Message()
			if s, ok := wktToString(msg); ok {
				// 与标量 list 保持一致：repeated 形式为 key=v0&key=v1...
				query.Add(key, s)
				continue
			}
			for k, v := range parseMessage(msg) {
				newKey := fmt.Sprintf("%s[%d].%s", key, i, k)
				query[newKey] = append(query[newKey], v...)
			}
		}
	}
	return query
}

func mapToValues(fd protoreflect.FieldDescriptor, m protoreflect.Map) url.Values {
	query := url.Values{}
	switch fd.MapValue().Kind() {
	case protoreflect.StringKind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			query.Add(getMapKey(fd.MapKey(), k.Value()), v.String())
			return true
		})
	case protoreflect.BoolKind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			query.Add(getMapKey(fd.MapKey(), k.Value()), strconv.FormatBool(v.Bool()))
			return true
		})
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			query.Add(getMapKey(fd.MapKey(), k.Value()), strconv.FormatInt(v.Int(), 10))
			return true
		})
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			query.Add(getMapKey(fd.MapKey(), k.Value()), strconv.FormatUint(v.Uint(), 10))
			return true
		})
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			query.Add(getMapKey(fd.MapKey(), k.Value()), strconv.FormatFloat(v.Float(), 'f', -1, 64))
			return true
		})
	case protoreflect.BytesKind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			query.Add(getMapKey(fd.MapKey(), k.Value()), base64.URLEncoding.EncodeToString(v.Bytes()))
			return true
		})
	case protoreflect.EnumKind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			query.Add(getMapKey(fd.MapKey(), k.Value()), enumToString(fd.MapValue().Enum(), v.Enum()))
			return true
		})
	case protoreflect.MessageKind, protoreflect.GroupKind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			mapKey := getMapKey(fd.MapKey(), k.Value())
			if s, ok := wktToString(v.Message()); ok {
				query.Add(mapKey, s)
				return true
			}
			for subK, subV := range parseMessage(v.Message()) {
				newKey := mapKey + "." + subK
				query[newKey] = append(query[newKey], subV...)
			}
			return true
		})
	}
	return query
}

// getMapKey 格式化 proto map 的 key。
// proto 规范规定 map key 只能是「整型/字符串/布尔型」
func getMapKey(fd protoreflect.FieldDescriptor, key protoreflect.Value) string {
	switch fd.Kind() {
	case protoreflect.StringKind:
		return key.String()
	case protoreflect.BoolKind:
		return strconv.FormatBool(key.Bool())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return strconv.FormatInt(key.Int(), 10)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return strconv.FormatUint(key.Uint(), 10)
	}
	return key.String()
}

func enumToString(ed protoreflect.EnumDescriptor, n protoreflect.EnumNumber) string {
	if ev := ed.Values().ByNumber(n); ev != nil {
		return string(ev.Name())
	}
	// 未知 enum number：输出数字，客户端可继续按数字解析
	return strconv.FormatInt(int64(n), 10)
}

// wktToString 把 Well-Known Types 序列化为扁平字符串，
// 对称于 getProtoMessage。返回 (s, true) 表示命中 WKT；
// 否则返回 ("", false)，调用方按普通 message 处理。
func wktToString(m protoreflect.Message) (string, bool) {
	switch msg := m.Interface().(type) {
	case *timestamppb.Timestamp:
		return msg.AsTime().UTC().Format(time.RFC3339Nano), true
	case *durationpb.Duration:
		return msg.AsDuration().String(), true
	case *wrapperspb.DoubleValue:
		return strconv.FormatFloat(msg.Value, 'f', -1, 64), true
	case *wrapperspb.FloatValue:
		return strconv.FormatFloat(float64(msg.Value), 'f', -1, 32), true
	case *wrapperspb.Int64Value:
		return strconv.FormatInt(msg.Value, 10), true
	case *wrapperspb.Int32Value:
		return strconv.FormatInt(int64(msg.Value), 10), true
	case *wrapperspb.UInt64Value:
		return strconv.FormatUint(msg.Value, 10), true
	case *wrapperspb.UInt32Value:
		return strconv.FormatUint(uint64(msg.Value), 10), true
	case *wrapperspb.BoolValue:
		return strconv.FormatBool(msg.Value), true
	case *wrapperspb.StringValue:
		return msg.Value, true
	case *wrapperspb.BytesValue:
		return base64.URLEncoding.EncodeToString(msg.Value), true
	case *fieldmaskpb.FieldMask:
		return strings.Join(msg.Paths, ","), true
	case *structpb.Value, *structpb.Struct, *structpb.ListValue:
		b, err := protojson.Marshal(msg)
		if err != nil {
			return err.Error(), true
		}
		return string(b), true
	}
	return "", false
}

// ValuesToMessage parses url.Values into proto.Message
func ValuesToMessage(msg proto.Message, values url.Values) error {
	return parseValues(msg.ProtoReflect(), valuesToMap(values))
}

func parseValues(msg protoreflect.Message, data map[string]any) error {
	fields := msg.Descriptor().Fields()
	for key, val := range data {
		fd := fields.ByName(protoreflect.Name(key))
		if fd == nil {
			continue
		}
		// 类型断言
		switch ss := val.(type) {
		case []string:
			// 不考虑map类型
			if len(ss) == 0 || fd.IsMap() {
				continue
			}
			// list
			if fd.IsList() {
				list := msg.Mutable(fd).List()
				for _, s := range ss {
					v, err := getProtoValue(fd, s)
					if err != nil {
						return err
					}
					if v.IsValid() {
						list.Append(v)
					}
				}
				continue
			}
			// 其他类型
			v, err := getProtoValue(fd, ss[0])
			if err != nil {
				return err
			}
			if v.IsValid() {
				msg.Set(fd, v)
			}
		case map[string]any:
			if fd.IsMap() {
				if err := setValuesMap(msg, fd, ss); err != nil {
					return err
				}
				continue
			}
			if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
				subMsg := msg.Mutable(fd).Message()
				if err := parseValues(subMsg, ss); err != nil {
					return err
				}
			}
		case []map[string]any:
			if fd.IsList() {
				list := msg.Mutable(fd).List()
				for _, m := range ss {
					subMsg := list.NewElement().Message()
					if err := parseValues(subMsg, m); err != nil {
						return err
					}
					list.Append(protoreflect.ValueOfMessage(subMsg))
				}
			}
		}
	}
	return nil
}

func getProtoValue(fd protoreflect.FieldDescriptor, s string) (protoreflect.Value, error) {
	var value protoreflect.Value
	switch fd.Kind() {
	case protoreflect.StringKind:
		value = protoreflect.ValueOfString(s)
	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(s)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfBool(v)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfInt32(int32(v))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfInt64(v)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfUint32(uint32(v))
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfUint64(v)
	case protoreflect.FloatKind:
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfFloat32(float32(v))
	case protoreflect.DoubleKind:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfFloat64(v)
	case protoreflect.BytesKind:
		b, err := Bytes(s)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfBytes(b)
	case protoreflect.EnumKind:
		enumVals := fd.Enum().Values()
		// 先按名字解析
		if ev := enumVals.ByName(protoreflect.Name(s)); ev != nil {
			value = protoreflect.ValueOfEnum(ev.Number())
			break
		}
		// 回退按数字解析（与 protojson 对齐）
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return value, fmt.Errorf("invalid enum value %q for %s", s, fd.Enum().FullName())
		}
		if ev := enumVals.ByNumber(protoreflect.EnumNumber(n)); ev != nil {
			value = protoreflect.ValueOfEnum(ev.Number())
			break
		}
		// proto3 允许未知 enum number，按数字透传
		value = protoreflect.ValueOfEnum(protoreflect.EnumNumber(n))
	case protoreflect.MessageKind, protoreflect.GroupKind:
		m, err := getProtoMessage(fd.Message(), s)
		if err != nil {
			return value, err
		}
		value = protoreflect.ValueOfMessage(m)
	}
	return value, nil
}

func parseTimestamp(s string) (*timestamppb.Timestamp, error) {
	t := carbon.Parse(s)
	if t.IsInvalid() {
		return nil, fmt.Errorf("cannot parse timestamp %q: unrecognized format", s)
	}

	v := timestamppb.New(t.StdTime())
	if !v.IsValid() {
		return nil, fmt.Errorf("timestamp %q is before 0001-01-01", s)
	}
	return v, nil
}

func getProtoMessage(md protoreflect.MessageDescriptor, s string) (protoreflect.Message, error) {
	var msg proto.Message
	switch md.FullName() {
	case "google.protobuf.Timestamp":
		v, err := parseTimestamp(s)
		if err != nil {
			return nil, err
		}
		msg = v
	case "google.protobuf.Duration":
		d, err := time.ParseDuration(s)
		if err != nil {
			return nil, err
		}
		msg = durationpb.New(d)
	case "google.protobuf.DoubleValue":
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		msg = wrapperspb.Double(v)
	case "google.protobuf.FloatValue":
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		msg = wrapperspb.Float(float32(v))
	case "google.protobuf.Int64Value":
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		msg = wrapperspb.Int64(v)
	case "google.protobuf.Int32Value":
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}
		msg = wrapperspb.Int32(int32(v))
	case "google.protobuf.UInt64Value":
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		msg = wrapperspb.UInt64(v)
	case "google.protobuf.UInt32Value":
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}
		msg = wrapperspb.UInt32(uint32(v))
	case "google.protobuf.BoolValue":
		v, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		msg = wrapperspb.Bool(v)
	case "google.protobuf.StringValue":
		msg = wrapperspb.String(s)
	case "google.protobuf.BytesValue":
		b, err := Bytes(s)
		if err != nil {
			return nil, err
		}
		msg = wrapperspb.Bytes(b)
	case "google.protobuf.FieldMask":
		msg = &fieldmaskpb.FieldMask{
			Paths: strings.Split(s, ","),
		}
	case "google.protobuf.Value":
		v := new(structpb.Value)
		if err := protojson.Unmarshal([]byte(s), v); err != nil {
			return nil, err
		}
		msg = v
	case "google.protobuf.Struct":
		v := new(structpb.Struct)
		if err := protojson.Unmarshal([]byte(s), v); err != nil {
			return nil, err
		}
		msg = v
	case "google.protobuf.ListValue":
		v := new(structpb.ListValue)
		if err := protojson.Unmarshal([]byte(s), v); err != nil {
			return nil, err
		}
		msg = v
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedMessage, string(md.FullName()))
	}
	return msg.ProtoReflect(), nil
}

func setValuesMap(msg protoreflect.Message, fd protoreflect.FieldDescriptor, data map[string]any) error {
	keyFD := fd.MapKey()
	valFD := fd.MapValue()

	m := msg.Mutable(fd).Map()
	for subKey, subVal := range data {
		keyVal, err := getProtoValue(keyFD, subKey)
		if err != nil {
			return err
		}
		if mapKey := keyVal.MapKey(); mapKey.IsValid() {
			switch v := subVal.(type) {
			case []string:
				if len(v) == 0 {
					continue
				}
				mapVal, err := getProtoValue(valFD, v[0])
				if err != nil {
					kind := valFD.Kind()
					if kind == protoreflect.MessageKind || kind == protoreflect.GroupKind {
						if errors.Is(err, ErrUnsupportedMessage) {
							continue
						}
					}
					return err
				}
				if mapVal.IsValid() {
					m.Set(mapKey, mapVal)
				}
			case map[string]any:
				// 仅支持message类型
				if valFD.Kind() == protoreflect.MessageKind || valFD.Kind() == protoreflect.GroupKind {
					subMsg := m.NewValue().Message()
					if err := parseValues(subMsg, v); err != nil {
						return err
					}
					m.Set(mapKey, protoreflect.ValueOfMessage(subMsg))
				}
			}
		}
	}
	return nil
}

// valuesToMap 将 url.Values 转换为嵌套 map，供 parseValues 使用。
// 支持点号嵌套（a.b.c）与下标（items[0].name）。
// 下标仅作为"分组键"使用，不决定位置与容量：
//
//	friends[-1]=a&friends[0]=b&friends[10]=c -> []map{...} 长度为 3，按下标升序。
//
// 非法下标（如 [abc]）、类型冲突 silent skip 该 key，不 panic。
func valuesToMap(values url.Values) map[string]any {
	root := make(map[string]any)
	for key, rawVal := range values {
		parts := strings.Split(flattenKey(key), ".")
		writeTree(root, parts, rawVal)
	}
	return finalizeArrays(root).(map[string]any)
}

func writeTree(root map[string]any, parts []string, value []string) {
	var cur any = root
	for i, part := range parts {
		isLast := i == len(parts)-1

		lb := strings.IndexByte(part, '[')
		if lb < 0 {
			parent, ok := cur.(map[string]any)
			if !ok {
				return
			}
			if isLast {
				parent[part] = value
				return
			}
			raw, exists := parent[part]
			if !exists {
				raw = make(map[string]any)
				parent[part] = raw
			}
			next, ok := raw.(map[string]any)
			if !ok {
				return
			}
			cur = next
			continue
		}

		rb := strings.IndexByte(part, ']')
		if rb <= lb {
			return // 畸形 [ 无匹配 ]
		}
		idx, err := strconv.Atoi(part[lb+1 : rb])
		if err != nil {
			return // 非整数下标（如 [abc]），跳过
		}
		arrayKey := part[:lb]

		parent, ok := cur.(map[string]any)
		if !ok {
			return
		}
		raw, exists := parent[arrayKey]
		if !exists {
			raw = make(map[int]any)
			parent[arrayKey] = raw
		}
		arr, ok := raw.(map[int]any)
		if !ok {
			return
		}
		if isLast {
			// 末段是 tags[i]=v 形式：允许标量 list
			switch t := arr[idx].(type) {
			case []string:
				arr[idx] = append(t, value...)
			case map[string]any:
				// 已有嵌套消息，忽略扁平覆盖，避免用户输入冲突导致数据丢失
				return
			default:
				arr[idx] = value
			}
			return
		}
		next, exists := arr[idx]
		if !exists {
			next = make(map[string]any)
			arr[idx] = next
		}
		if _, ok := next.(map[string]any); !ok {
			return
		}
		cur = next
	}
}

// finalizeArrays 递归地把 map[int]any 节点按 key 升序展开为紧致切片：
//   - 若桶内全是 []string，则合并为 []string（标量 list）；
//   - 否则为 []map[string]any（message list）。
//
// 下标只决定顺序，不影响结果长度；稀疏下标会被压紧。
func finalizeArrays(v any) any {
	switch t := v.(type) {
	case map[string]any:
		for k, sub := range t {
			t[k] = finalizeArrays(sub)
		}
		return t
	case map[int]any:
		keys := make([]int, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		scalar := true
		for _, k := range keys {
			if _, ok := t[k].(map[string]any); ok {
				scalar = false
				break
			}
		}
		if scalar {
			ss := make([]string, 0, len(keys))
			for _, k := range keys {
				if s, ok := t[k].([]string); ok {
					ss = append(ss, s...)
				}
			}
			return ss
		}
		items := make([]map[string]any, 0, len(keys))
		for _, k := range keys {
			if m, ok := t[k].(map[string]any); ok {
				items = append(items, finalizeArrays(m).(map[string]any))
			}
		}
		return items
	default:
		return v
	}
}

func flattenKey(key string) string {
	items := strings.FieldsFunc(key, func(r rune) bool {
		return r == '[' || r == ']'
	})
	if len(items) <= 1 {
		return key
	}

	var builder strings.Builder
	builder.WriteString(items[0])
	for _, v := range items[1:] {
		if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			builder.WriteString("[")
			builder.WriteString(v)
			builder.WriteString("]")
		} else {
			if !strings.HasPrefix(v, ".") {
				builder.WriteString(".")
			}
			builder.WriteString(v)
		}
	}
	return builder.String()
}
