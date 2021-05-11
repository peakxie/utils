package slicex

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gogo/protobuf/sortkeys"
)

// 去除Slice中重复值
func RemoveDuplicates(s interface{}) interface{} {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	switch typ.Kind() {
	case reflect.Slice:
		if val.Len() == 0 {
			return s
		}
		m := make(map[interface{}]int)
		for i := 0; i < val.Len(); i++ {
			m[val.Index(i).Interface()]++
		}
		res := reflect.MakeSlice(typ, 0, 0)
		for k := range m {
			res = reflect.Append(res, reflect.ValueOf(k))
		}
		return res.Interface()
	default:
		return s
	}
}

// 检查Slice, a是否包含b
func CheckSubset(a, b []uint64) bool {
	if len(b) == 0 {
		return true
	}
	t := RemoveDuplicates(a[0:]).([]uint64)
	s := RemoveDuplicates(b[0:]).([]uint64)
	tLen := len(t)
	sLen := len(s)
	if tLen < sLen {
		return false
	}
	//sort.Sort(Uint64Slice(t))
	//sort.Sort(Uint64Slice(s))
	sortkeys.Uint64s(t)
	sortkeys.Uint64s(s)
	for i, j := 0, 0; i+sLen <= tLen; i++ {
		if sLen == 0 {
			return true
		}
		if t[i] == s[j] {
			j++
			sLen--
		}
	}
	return false
}

//两个Slice的差集，即a-b的结果
func DifferenceSet(a, b interface{}) interface{} {
	valA := reflect.ValueOf(a)
	typA := reflect.TypeOf(a)

	valB := reflect.ValueOf(b)
	typB := reflect.TypeOf(b)

	if typA.Kind() != reflect.Slice || typB.Kind() != reflect.Slice {
		return a
	}

	if valA.Len() == 0 || valB.Len() == 0 {
		return a
	}

	c := reflect.MakeSlice(typA, 0, 0)
	for i := 0; i < valA.Len(); i++ {
		exist := false
		for j := 0; j < valB.Len(); j++ {
			if reflect.DeepEqual(valA.Index(i).Interface(), valB.Index(j).Interface()) {
				exist = true
				break
			}
		}
		if exist == false {
			c = reflect.Append(c, valA.Index(i))
		}
	}

	return c.Interface()
}

// 求json数组交集
func JsonOverlaps(jsonArray string, filterArray interface{}) string {
	var condition string
	val := reflect.ValueOf(filterArray)
	typ := reflect.TypeOf(filterArray)
	switch typ.Kind() {
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if i != 0 {
				condition += " OR "
			}
			condition += fmt.Sprintf("JSON_CONTAINS(%s, JSON_ARRAY(%v))", jsonArray, val.Index(i))
		}

	default:
		return ""
	}
	return condition
}

// 求json数组并
func UnJsonOverlaps(jsonArray string, filterArray interface{}) string {
	var condition string
	val := reflect.ValueOf(filterArray)
	typ := reflect.TypeOf(filterArray)
	switch typ.Kind() {
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if i != 0 {
				condition += " AND "
			}
			condition += fmt.Sprintf("!JSON_CONTAINS(%s, JSON_ARRAY(%v))", jsonArray, val.Index(i))
		}

	default:
		return ""
	}
	return condition
}

func Uint64SliceToStringSlice(in []uint64) []string {
	var out []string
	for _, v := range in {
		item := strconv.FormatUint(v, 10)
		out = append(out, item)
	}
	return out
}

func StringSliceToUint64Slice(in []string) []uint64 {
	var out []uint64
	for _, v := range in {
		u, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			continue
		}
		out = append(out, u)
	}
	return out
}
