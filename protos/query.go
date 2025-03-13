package protos

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// MessageToQuery parses proto.Message into url.Values
func MessageToQuery(msg proto.Message) url.Values {
	return parseMessage(msg.ProtoReflect())
}

func parseMessage(msg protoreflect.Message) url.Values {
	query := url.Values{}
	fields := msg.Descriptor().Fields()

	for index, count := 0, 0; count < fields.Len(); index++ {
		fd := fields.ByNumber(protoreflect.FieldNumber(index))
		if fd == nil {
			continue
		}
		count++
		key := string(fd.Name())
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
			query.Add(string(fd.Name()), base64.URLEncoding.EncodeToString(val.Bytes()))
		case protoreflect.EnumKind:
			e := fd.Enum().Values()
			num := val.Enum()
			query.Add(key, string(e.ByNumber(num).Name()))
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
			e := fd.Enum().Values()
			num := list.Get(i).Enum()
			query.Add(key, string(e.ByNumber(num).Name()))
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
			e := fd.MapValue().Enum().Values()
			query.Add(getMapKey(fd.MapKey(), k.Value()), string(e.ByNumber(v.Enum()).Name()))
			return true
		})
	case protoreflect.MessageKind, protoreflect.GroupKind:
		m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			key := getMapKey(fd.MapKey(), k.Value())
			for k, v := range parseMessage(v.Message()) {
				newKey := key + "." + k
				query[newKey] = append(query[newKey], v...)
			}
			return true
		})
	}
	return query
}

func getMapKey(fd protoreflect.FieldDescriptor, key protoreflect.Value) string {
	switch fd.Kind() {
	case protoreflect.StringKind:
		return key.String()
	case protoreflect.BoolKind:
		return strconv.FormatBool(key.Bool())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return strconv.FormatInt(key.Int(), 10)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return strconv.FormatUint(key.Uint(), 10)
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return strconv.FormatFloat(key.Float(), 'f', -1, 64)
	case protoreflect.BytesKind:
		return base64.URLEncoding.EncodeToString(key.Bytes())
	case protoreflect.EnumKind:
		num := key.Enum()
		return string(fd.Enum().Values().ByNumber(num).Name())
	}
	return key.String()
}

// ValuesToMessage parses url.Values into proto.Message
func ValuesToMessage(msg proto.Message, values url.Values) error {
	return parseValues(msg.ProtoReflect(), valuesToMap(values))
}

func parseValues(msg protoreflect.Message, query map[string]any) error {
	fields := msg.Descriptor().Fields()
	for key, val := range query {
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

func getProtoMessage(md protoreflect.MessageDescriptor, s string) (protoreflect.Message, error) {
	var msg proto.Message
	switch md.FullName() {
	case "google.protobuf.Timestamp":
		t, err := time.Parse(time.DateTime, s)
		if err != nil {
			return nil, err
		}
		v := timestamppb.New(t)
		if ok := v.IsValid(); !ok {
			return nil, fmt.Errorf("%s before 0001-01-01", s)
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
		v := new(fieldmaskpb.FieldMask)
		v.Paths = append(v.Paths, strings.Split(s, ",")...)
		msg = v
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
	m := msg.Mutable(fd).Map()

	keyFD := fd.MapKey()
	valFD := fd.MapValue()

	for subKey, subVal := range data {
		keyVal, err := getProtoValue(keyFD, subKey)
		if err != nil {
			return err
		}
		if mapKey := keyVal.MapKey(); mapKey.IsValid() {
			switch v := subVal.(type) {
			case []string:
				// 不考虑message类型
				if len(v) != 0 && valFD.Kind() != protoreflect.MessageKind && valFD.Kind() != protoreflect.GroupKind {
					mapVal, err := getProtoValue(valFD, v[0])
					if err != nil {
						return err
					}
					if mapVal.IsValid() {
						m.Set(mapKey, mapVal)
					}
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

func valuesToMap(values url.Values) map[string]any {
	result := make(map[string]any)
	for key, value := range values {
		items := strings.Split(flattenKey(key), ".")
		if len(items) < 2 {
			result[key] = value
			continue
		}
		nestedMap := result
		for i, k := range items {
			if strings.Contains(k, "[") {
				arrayKey := k[:strings.Index(k, "[")]
				if _, ok := nestedMap[arrayKey]; !ok {
					nestedMap[arrayKey] = make([]map[string]any, 0)
				}
				array := nestedMap[arrayKey].([]map[string]any)
				index, _ := strconv.Atoi(k[strings.Index(k, "[")+1 : strings.Index(k, "]")])
				for len(array) <= index {
					array = append(array, make(map[string]any))
				}
				nestedMap[arrayKey] = array
				nestedMap = array[index]
			} else {
				if i == len(items)-1 {
					nestedMap[k] = value
				} else {
					if _, ok := nestedMap[k]; !ok {
						nestedMap[k] = make(map[string]any)
					}
					nestedMap = nestedMap[k].(map[string]any)
				}
			}
		}
	}
	return result
}

func flattenKey(key string) string {
	items := extractBracketKey(key)
	if len(items) == 1 {
		return key
	}
	var builder strings.Builder
	builder.WriteString(items[0])
	for _, v := range items[1:] {
		if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			builder.WriteString("[" + v + "]")
		} else {
			if v[0] != '.' {
				builder.WriteString(".")
			}
			builder.WriteString(v)
		}
	}
	return builder.String()
}

func extractBracketKey(key string) []string {
	reg := regexp.MustCompile(`\[|\]`)
	parts := reg.Split(key, -1)
	var result []string
	for _, v := range parts {
		if len(v) != 0 {
			result = append(result, v)
		}
	}
	return result
}
