package loghelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
)

//copy from genejli's util

//add by gene
//打印请求包或者回包时，对报文内容进行处理，主要针对pb：
// 1、如果pb协议没有过长，直接使用pb.String()
// 2、过长则对结构体中的字段进行处理，参考下面的可配置变量

var (
	// PrintProtoLen 设置打印pb报文的长度，超过这个长度将会对报文中的超长字段进行处理（截断）
	PrintProtoLen = 512

	// PrintStringLen 设置打印报文中字符串的长度，超过这个长度将会对报文中的超长字段进行处理（截断）
	PrintStringLen = 64

	// PrintSliceLen 设置打印报文中的slice或者array的长度，超过这个长度将会对报文中的超长内容进行处理（截断）
	PrintSliceLen = 64

	// 处理string的方式，默认截断后面补上 "..."
	DealStringHook DealStringHookFunc = getStringSlice
)

type DealStringHookFunc func(string) string

// 把定义了Json字段的结构本转换为Json字符串输出
func ToJsonString(r interface{}) string {
	b, _ := json.Marshal(r)
	return string(b)
}

type SelfJsonString interface {
	ToJson() string
}

// String 把i转换成string
func String(i interface{}) string {
	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.String {
		return reflect.ValueOf(i).String()
	}
	if typ.Kind() == reflect.Ptr {
		return String(reflect.ValueOf(i).Elem())
	}
	return fmt.Sprintf("%v", i)
}

// ProtoToPrintString pb协议转换成可打印的json string，会对超长内容截断
func ProtoToPrintString(p proto.Message) string {
	if p == nil {
		return ""
	}
	str := p.String()
	if PrintProtoLen <= 0 || len(str) < PrintProtoLen {
		return str
	}
	m, e := GetStructFields(p)
	if e == nil {
		data, e := json.Marshal(m)
		if e == nil {
			return string(data)
		}
	}
	return ""

}

// StructToPrintString 结构体转换成可打印的json string，会对超长内容截断
func StructToPrintString(p interface{}) string {
	if p == nil {
		return ""
	}
	if pm, ok := p.(proto.Message); ok {
		return ProtoToPrintString(pm)
	}
	m, e := GetStructFields(p)
	if e == nil {
		data, e := json.Marshal(m)
		if e == nil {
			return string(data)
		}
	}
	return ""
}

// ToPrintString 入参换成可打印的json string，包括protobuf和普通结构体，或者其他类型
func ToPrintString(i interface{}) string {
	if i == nil {
		return ""
	}
	typ := reflect.TypeOf(i)
	kind := typ.Kind()
	if typ.Kind() == reflect.Ptr {
		if reflect.ValueOf(i).IsNil() {
			return "nil"
		}
		kind = typ.Elem().Kind()
	}
	switch kind {
	case reflect.Invalid:
		return "invalid"
	case reflect.Struct:
		return StructToPrintString(i)
	default:
		return ToJsonString(i)
	}
}

func hasJsonName(field reflect.StructField) bool {
	tag := field.Tag.Get("json")
	if tag != "" {
		splitTags := strings.Split(tag, ",")
		if splitTags[0] != "" && splitTags[0] != "-" {
			return true
		}
	}
	return false
}

func getJsonName(field reflect.StructField) string {
	tagName := field.Name
	tag := field.Tag.Get("json")
	if tag != "" {
		splitTags := strings.Split(tag, ",")
		if splitTags[0] != "" {
			tagName = splitTags[0]
		}
	}
	return tagName
}

//
//func handleString(m map[string]interface{}, key string, s interface{}) {
//	if s, ok := s.(string); ok {
//		if len(s) < PrintStringLen {
//			m[key] = s
//		} else {
//			m[key] = s[:PrintStringLen] + "..."
//		}
//		return
//	}
//
//	if s, ok := s.(*string); ok {
//		if len(*s) < PrintStringLen {
//			m[key] = s
//		} else {
//			m[key] = (*s)[:PrintStringLen] + "..."
//		}
//		return
//	}
//}

func getStringSlice(s string) string {
	if PrintStringLen <= 0 || len(s) < PrintStringLen {
		return s
	} else {
		return s[:PrintStringLen] + "..."
	}
}

func handlerField(fields map[string]interface{}, fieldName string, field reflect.Value, t reflect.Type, anonymous bool) {
	// type field
	if _, ok := fields[fieldName]; ok {
		return
	}
	kind := t.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		var values []interface{}
		moreFields := false
		fLen := field.Len()
		if PrintSliceLen > 0 && fLen > PrintSliceLen {
			field = field.Slice(0, PrintSliceLen)
			moreFields = true
		}
		defer func() {
			if len(values) > 0 {
				if moreFields {
					//多余的字段省略，使用...代替
					values = append(values, fmt.Sprintf("(more %d item)...", fLen-PrintSliceLen))
				}
				fields[fieldName] = values
			}
		}()
		var sliceKind reflect.Kind
		bPtr := false
		if field.Type().Elem().Kind() == reflect.Ptr {
			sliceKind = field.Type().Elem().Elem().Kind()
			bPtr = true
		} else {
			sliceKind = field.Type().Elem().Kind()
		}
		switch sliceKind {
		case reflect.Bool:
			fallthrough
		case reflect.Float32, reflect.Float64:
			fallthrough
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			//fields[fieldName] = field.Interface()
			for i := 0; i < field.Len(); i++ {
				values = append(values, field.Index(i).Interface())
			}

		case reflect.String:
			for i := 0; i < field.Len(); i++ {
				var t *string
				if bPtr {
					if !field.Index(i).IsNil() {
						str := DealStringHook(field.Index(i).Elem().String())
						t = &str
					}
				} else {
					str := DealStringHook(field.Index(i).String())
					t = &str
				}
				values = append(values, t)
			}
		case reflect.Struct:
			for i := 0; i < field.Len(); i++ {
				t := make(map[string]interface{})
				if field.Type().Elem().Kind() == reflect.Ptr {
					if !field.Index(i).IsNil() {
						handlerStruct(t, reflect.ValueOf(field.Index(i).Elem().Interface()), reflect.TypeOf(field.Index(i).Elem().Interface()))
					}
				} else {
					handlerStruct(t, reflect.ValueOf(field.Index(i).Interface()), reflect.TypeOf(field.Index(i).Interface()))
				}
				if len(t) > 0 {
					values = append(values, t)
				}
			}
		case reflect.Interface:
			for i := 0; i < field.Len(); i++ {
				t := make(map[string]interface{})
				if !field.Index(i).IsNil() {
					handlerField(t, fieldName, field.Index(i).Elem(), field.Index(i).Elem().Type(), false)
				}
				if len(t) > 0 {
					values = append(values, t[fieldName])
				}
			}
		case reflect.Map:
			for i := 0; i < field.Len(); i++ {
				t := make(map[string]interface{})
				if field.Type().Elem().Kind() == reflect.Ptr {
					if !field.Index(i).IsNil() {
						handlerField(t, fieldName, field.Index(i).Elem(), field.Index(i).Elem().Type(), false)
					}
				} else {
					handlerField(t, fieldName, field.Index(i), field.Index(i).Type(), false)
				}
				if len(t) > 0 {
					values = append(values, t[fieldName])
				}
			}
		default:
			fmt.Printf("type:%v", sliceKind)
			//panic("reflect.TypeOf(param).Elem().Kind() no setting")
		}
	case reflect.String:
		str := DealStringHook(field.String())
		fields[fieldName] = str
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Bool:
		fields[fieldName] = field.Interface()
	case reflect.Struct:
		if anonymous {
			handlerStruct(fields, reflect.ValueOf(field.Interface()), reflect.TypeOf(field.Interface()))
		} else {
			temp := make(map[string]interface{})
			handlerStruct(temp, reflect.ValueOf(field.Interface()), reflect.TypeOf(field.Interface()))
			fields[fieldName] = temp
		}
	case reflect.Ptr:
		if !field.IsNil() {
			handlerField(fields, fieldName, field.Elem(), t.Elem(), false)
		}
	case reflect.Interface:
		if !field.IsNil() {
			handlerField(fields, fieldName, field.Elem(), field.Elem().Type(), false)
		}
	case reflect.Map:
		if !field.IsNil() {
			temp := make(map[string]interface{})
			for _, key := range field.MapKeys() {
				handlerField(temp, formatAtom(key), field.MapIndex(key), field.MapIndex(key).Type(), false)
			}
			fields[fieldName] = temp
		}
	default:
		//panic("reflect.Type.Kind() no setting")
	}
	return
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return v.String()
	case reflect.Chan, reflect.Func, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Interface, reflect.Ptr:
		return formatAtom(v.Elem())
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}
func handlerStruct(fields map[string]interface{}, v reflect.Value, t reflect.Type) {
	for i := 0; i < v.Type().NumField(); i++ {
		if !v.Field(i).CanInterface() || !v.Field(i).IsValid() || (isOmitEmpty(v.Type().Field(i)) && isEmptyValue(v.Field(i))) {
			continue
		}
		fieldName := getJsonName(t.Field(i))
		if fieldName != "-" {
			handlerField(fields, fieldName, v.Field(i), v.Type().Field(i).Type, t.Field(i).Anonymous)
		}
	}
}

func isOmitEmpty(field reflect.StructField) bool {
	tag := field.Tag.Get("json")
	if strings.Contains(tag, "omitempty") {
		return true
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// GetStructFields 把结构体转换成一个map
func GetStructFields(st interface{}) (fields map[string]interface{}, err error) {
	fields = make(map[string]interface{})
	if st == nil {
		return
	}
	v := reflect.ValueOf(st)
	t := reflect.TypeOf(st)

	if v.Kind() == reflect.Ptr {
		if v.Elem().Type().Kind() == reflect.Struct {
			handlerStruct(fields, v.Elem(), v.Elem().Type())
			return
		}
	}
	switch v.Kind() {
	case reflect.Struct:
		handlerStruct(fields, v, t)
	default:
		err = errors.New("Can't handler type: " + v.Type().String())
	}
	return
}
