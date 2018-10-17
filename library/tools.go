package library

import (
	"strconv"
	"fmt"
	"reflect"
)

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func CheckQueryError(err error, query string)  {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s %s", query, err))
	}
}

func  ToString(arg interface{}) string {
	switch arg.(type) {
	case int:
		return strconv.Itoa(arg.(int))
	case int64:
		return strconv.FormatInt(arg.(int64), 10)
	case []byte:
		return string(arg.([]byte))
	case string:
		return arg.(string)
	default:
		return ""
	}
}

func Addslashes(v string) string {
	pos := 0
	buf := make([]byte, len(v)*2)
	for i := 0; i < len(v); i++ {
		c := v[i]
		if c == '\'' || c == '"' || c == '\\' {
			buf[pos] = '\\'
			buf[pos+1] = c
			pos += 2
		} else {
			buf[pos] = c
			pos++
		}
	}
	return string(buf[:pos])
}


func ToInt64(value interface{}) int64 {
	val := reflect.ValueOf(value)
	var d int64
	var err error
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	case string:
		d, err = strconv.ParseInt(val.String(), 10, 64)
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	if err != nil {
		d = 0
	}
	return d
}

func ToInt(value interface{}) int {
	val := reflect.ValueOf(value)
	var d int
	var err error
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = int(val.Int())
	case uint, uint8, uint16, uint32, uint64:
		d = int(val.Uint())
	case string:
		d, err = strconv.Atoi(val.String())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	if err != nil {
		d = 0
	}

	return d
}