package protokit

import (
	"encoding/base64"
	"fmt"
	"net/url"
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

// MessageToValues parses proto.Message into url.Values.
// Only fields that are explicitly set (non-default) are included.
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

		if fd.IsList() {
			for k, v := range listToValues(fd, val.List()) {
				query[k] = append(query[k], v...)
			}
			continue
		}

		if fd.IsMap() {
			for k, v := range mapToValues(fd, val.Map()) {
				newKey := key + "." + k
				query[newKey] = append(query[newKey], v...)
			}
			continue
		}

		switch fd.Kind() {
		case protoreflect.StringKind:
			query.Add(key, val.String())
		case protoreflect.BoolKind:
			query.Add(key, strconv.FormatBool(val.Bool()))
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			query.Add(key, strconv.FormatInt(val.Int(), 10))
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			query.Add(key, strconv.FormatUint(val.Uint(), 10))
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			query.Add(key, strconv.FormatFloat(val.Float(), 'f', -1, 64))
		case protoreflect.BytesKind:
			query.Add(key, base64.URLEncoding.EncodeToString(val.Bytes()))
		case protoreflect.EnumKind:
			e := fd.Enum().Values()
			query.Add(key, string(e.ByNumber(val.Enum()).Name()))
		case protoreflect.MessageKind, protoreflect.GroupKind:
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
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		for i := 0; i < list.Len(); i++ {
			query.Add(key, strconv.FormatInt(list.Get(i).Int(), 10))
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
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
			e := fd.Enum().Values()
			query.Add(key, string(e.ByNumber(list.Get(i).Enum()).Name()))
		}
	case protoreflect.MessageKind, protoreflect.GroupKind:
		for i := 0; i < list.Len(); i++ {
			for k, v := range parseMessage(list.Get(i).Message()) {
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
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			query.Add(getMapKey(fd.MapKey(), k.Value()), strconv.FormatInt(v.Int(), 10))
			return true
		})
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
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
			e := fd.MapValue().Enum().Values()
			query.Add(getMapKey(fd.MapKey(), k.Value()), string(e.ByNumber(v.Enum()).Name()))
			return true
		})
	case protoreflect.MessageKind, protoreflect.GroupKind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			mapKey := getMapKey(fd.MapKey(), k.Value())
			for subK, subV := range parseMessage(v.Message()) {
				newKey := mapKey + "." + subK
				query[newKey] = append(query[newKey], subV...)
			}
			return true
		})
	}
	return query
}

// proto 规范中 map key 只能是整数或字符串，不会走到这里
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
		switch ss := val.(type) {
		case []string:
			if len(ss) == 0 || fd.IsMap() {
				continue
			}
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
		enumDesc := fd.Enum()
		enumVal := enumDesc.Values().ByName(protoreflect.Name(s))
		if enumVal != nil {
			value = protoreflect.ValueOfEnum(enumVal.Number())
		}
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
	default:
		return nil, fmt.Errorf("unsupported message type: %s", string(md.FullName()))
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
		mapKey := keyVal.MapKey()
		if !mapKey.IsValid() {
			continue
		}
		switch v := subVal.(type) {
		case []string:
			if len(v) == 0 {
				continue
			}
			if valFD.Kind() == protoreflect.MessageKind || valFD.Kind() == protoreflect.GroupKind {
				continue
			}
			mapVal, err := getProtoValue(valFD, v[0])
			if err != nil {
				return err
			}
			if mapVal.IsValid() {
				m.Set(mapKey, mapVal)
			}
		case map[string]any:
			if valFD.Kind() != protoreflect.MessageKind && valFD.Kind() != protoreflect.GroupKind {
				continue
			}
			subMsg := m.NewValue().Message()
			if err := parseValues(subMsg, v); err != nil {
				return err
			}
			m.Set(mapKey, protoreflect.ValueOfMessage(subMsg))
		}
	}
	return nil
}

// valuesToMap converts url.Values into a nested map for parseValues.
// Supports dot-notation for nested fields and indexed brackets for repeated fields,
// e.g. "address.city", "items[0].name".
func valuesToMap(values url.Values) map[string]any {
	result := make(map[string]any)
	for key, value := range values {
		segments := splitKey(key)
		if len(segments) < 2 {
			result[key] = value
			continue
		}
		if err := setNested(result, segments, value); err != nil {
			// 类型冲突时跳过该 key，不 panic
			continue
		}
	}
	return result
}

// segment represents one parsed path component.
type segment struct {
	key   string
	index int  // >= 0 表示数组下标
	isArr bool // 是否为数组访问
}

// splitKey parses a dotted/bracketed key into segments.
// e.g. "items[0].name" -> [{key:"items",isArr:true,index:0}, {key:"name"}]
func splitKey(key string) []segment {
	var segments []segment
	parts := strings.Split(key, ".")
	for _, part := range parts {
		if part == "" {
			continue
		}
		lb := strings.Index(part, "[")
		rb := strings.Index(part, "]")
		if lb >= 0 && rb > lb {
			arrKey := part[:lb]
			idx, err := strconv.Atoi(part[lb+1 : rb])
			if err != nil {
				// 非数字下标，当普通 key 处理
				segments = append(segments, segment{key: part})
				continue
			}
			segments = append(segments, segment{key: arrKey, isArr: true, index: idx})
			// 处理 ] 后面还有内容的情况，如 "[0]suffix"（罕见但防御）
			if rb+1 < len(part) {
				segments = append(segments, segment{key: part[rb+1:]})
			}
		} else {
			segments = append(segments, segment{key: part})
		}
	}
	return segments
}

// setNested writes value into nested map/slice structure described by segments.
// Returns error on type conflict to allow caller to skip gracefully.
func setNested(root map[string]any, segs []segment, value []string) error {
	cur := root
	for i, seg := range segs {
		isLast := i == len(segs)-1

		if seg.isArr {
			// 当前层是数组
			raw, exists := cur[seg.key]
			if !exists {
				cur[seg.key] = make([]map[string]any, 0)
				raw = cur[seg.key]
			}

			arr, ok := raw.([]map[string]any)
			if !ok {
				return fmt.Errorf("type conflict at key %q: expected []map, got %T", seg.key, raw)
			}
			for len(arr) <= seg.index {
				arr = append(arr, make(map[string]any))
			}
			cur[seg.key] = arr

			if isLast {
				// 数组最后一层直接存 value（极少见，防御性处理）
				arr[seg.index][seg.key] = value
			} else {
				cur = arr[seg.index]
			}
		} else {
			if isLast {
				cur[seg.key] = value
			} else {
				raw, exists := cur[seg.key]
				if !exists {
					cur[seg.key] = make(map[string]any)
					raw = cur[seg.key]
				}

				next, ok := raw.(map[string]any)
				if !ok {
					return fmt.Errorf("type conflict at key %q: expected map, got %T", seg.key, raw)
				}
				cur = next
			}
		}
	}
	return nil
}
