// Package pocket @author K·J Create at 2019-04-09 15:10
package pocket

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	uuid "github.com/satori/go.uuid"
)

var (
	mutex sync.Mutex
)

// GetUUID gene uuid
func GetUUID() uuid.UUID {
	mutex.Lock()
	defer mutex.Unlock()
	return uuid.NewV4()
}

// UnixSecond unix time second
func UnixSecond() int64 {
	return time.Now().Unix()
}

// UnixMillisecond unix time millisecond
func UnixMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// pow 次方
func pow(x, n int) int32 {
	ret := 1 // 结果初始为0次方的值，整数0次方为1。如果是矩阵，则为单元矩阵。
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return int32(ret)
}

// CreateCaptcha random
func CreateCaptcha(digit int) string {
	format := "%0" + fmt.Sprintf("%d", digit) + "v"
	return fmt.Sprintf(format, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(pow(10, digit)))
}

// md5 md5 32 lowercase
func Md5(txt string) string {
	md5Hash := md5.New()
	io.WriteString(md5Hash, txt)
	md5Bytes := md5Hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(md5Bytes))
}

// Diff compile 2objects
func Diff(origin interface{}, current interface{}) (string, string) {
	originStr := ""
	currentStr := ""
	originType := reflect.TypeOf(origin)
	currentType := reflect.TypeOf(current)

	if originType.Kind() == reflect.Ptr {
		originType = originType.Elem()
	}
	if currentType.Kind() == reflect.Ptr {
		currentType = currentType.Elem()
	}

	if originType.Name() != currentType.Name() {
		return originStr, currentStr
	}

	if originType.Kind() != reflect.Struct || currentType.Kind() != reflect.Struct {
		return originStr, currentStr
	}
	originValue := reflect.ValueOf(origin)
	currentValue := reflect.ValueOf(current)

	for i := 0; i < originType.NumField(); i++ {
		name := originType.Field(i).Name
		if "UpdatedAt" == name {
			continue
		}
		var v1, v2 string
		if originValue.Field(i).Kind() == reflect.Ptr {
			v1 = getValue(originValue.Field(i).Elem())
		} else {
			v1 = getValue(originValue.Field(i))
		}
		if currentValue.Field(i).Kind() == reflect.Ptr {
			v2 = getValue(currentValue.Field(i).Elem())
		} else {
			v2 = getValue(currentValue.Field(i))
		}
		if ("ID" == name || "Id" == name) || ("" != v2 && v1 != v2) {
			if "" != originStr {
				originStr += ","
				currentStr += ","
			}

			originStr = originStr + name + ":" + v1
			currentStr = currentStr + name + ":" + v2
		}
	}
	return originStr, currentStr
}

func value(v reflect.Value) (string, error) {
	if v.CanInterface() {
		switch v.Kind() {
		case reflect.Ptr:
			if !v.IsNil() {
				return fmt.Sprintf("%+v", v.Elem().Interface()), nil
			}
		default:
			return fmt.Sprintf("%v", v.Interface()), nil
		}
	} else {
		if v.CanAddr() {
			addr := v.Addr()
			ptr := unsafe.Pointer(addr.Pointer())
			switch v.Type().String() {
			case "int":
				fallthrough
			case "*int":
				fallthrough
			case "int64":
				fallthrough
			case "*int64":
				return fmt.Sprintf("%d", *(*int)(ptr)), nil
			case "*string":
				fallthrough
			case "string":
				return *(*string)(ptr), nil
			default:
				return "", errors.New("not support type")
			}
		}
	}
	return "", nil
}

func ToString(origin interface{}, tag, joiner, separator string) string {
	originStr := ""
	originType := reflect.TypeOf(origin)
	if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
		return originStr
	}

	originValue := reflect.ValueOf(origin)

	for i := 0; i < originType.Elem().NumField(); i++ {
		if !("" == originType.Elem().Field(i).Tag.Get(tag)) {
			if originValue.Elem().Field(i).Kind() == reflect.Ptr {
				if !originValue.Elem().Field(i).IsNil() {
					if "" != originStr {
						originStr += separator
					}
					value := fmt.Sprintf("%+v", originValue.Elem().Field(i).Elem())
					originStr = originStr + originType.Elem().Field(i).Name + joiner + value
				}
			} else {
				if "" != originStr {
					originStr += separator
				}
				value := fmt.Sprint(originValue.Elem().Field(i).Interface())
				originStr = originStr + originType.Elem().Field(i).Name + joiner + value
			}

		}
	}
	return originStr
}

// DecodeQuery decode get query
func DecodeQuery(dst interface{}, src map[string][]string) error {
	t := reflect.TypeOf(dst)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("schema: interface must be a pointer to struct")
	}
	v := reflect.ValueOf(dst).Elem()
	for i := 0; i < t.Elem().NumField(); i++ {
		tag := t.Elem().Field(i).Tag.Get("schema")
		if "" != tag {
			if param, ok := src[tag]; ok {
				if len(param) > 0 && len(param[0]) > 0 {
					switch v.Field(i).Kind() {
					case reflect.Ptr:
						switch v.Field(i).Type().String() {
						case "*int":
							val, err := strconv.Atoi(param[0])
							if nil != err {
								Logger.Error(err.Error())
								break
							}
							p := new(int)
							*p = val
							v.Field(i).Set(reflect.ValueOf(p))
						case "*int64":
							val, err := strconv.ParseInt(param[0], 10, 64)
							if nil != err {
								Logger.Error(err.Error())
								break
							}
							p := new(int64)
							*p = val
							v.Field(i).Set(reflect.ValueOf(p))
						case "*string":
							p := new(string)
							*p = param[0]
							v.Field(i).Set(reflect.ValueOf(p))
						default:
							return errors.New(fmt.Sprintf("Not Support %s", tag))
						}
					case reflect.String:
						v.Field(i).SetString(param[0])
					case reflect.Int:
						val, err := strconv.ParseInt(param[0], 10, 32)
						if nil != err {
							Logger.Error(err.Error())
							break
						}
						v.Field(i).SetInt(val)
					case reflect.Int64:
						val, err := strconv.ParseInt(param[0], 10, 64)
						if nil != err {
							Logger.Error(err.Error())
							break
						}
						v.Field(i).SetInt(val)
					default:
						return errors.New(fmt.Sprintf("Not Support %s", tag))
					}
				}
			}
		}
	}
	return nil
}

func String(model interface{}) string {
	str := ""
	originType := reflect.TypeOf(model)
	if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
		Logger.Warn("param error")
		return ""
	}
	originValue := reflect.ValueOf(model)

	for i := 0; i < originType.Elem().NumField(); i++ {
		if originValue.Elem().Field(i).Kind() == reflect.Ptr {
			if !originValue.Elem().Field(i).IsNil() {
				if len(str) > 0 {
					str += ","
				}
				str += fmt.Sprintf("%s=%s", originType.Elem().Field(i).Name, getValue(originValue.Elem().Field(i).Elem()))
			}
		} else {
			val := getValue(originValue.Elem().Field(i))
			if "" == val {
				continue
			}
			if len(str) > 0 {
				str += ","
			}
			str += fmt.Sprintf("%s=%s", originType.Elem().Field(i).Name, val)
		}
	}
	return str
}

func getValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Slice, reflect.Array:
		arrayLen := v.Len()
		s := ""
		for i := 0; i < arrayLen; i++ {
			if i > 0 {
				s += ","
			}
			s += getValue(v.Index(i))
		}
		return s
	default:
		return ""
	}
}

// SnakeString, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func NilCondition(i interface{}) bool {
	return i == nil
}

func Filter(condition func(int, interface{}) bool, data interface{}) (interface{}, error) {
	dataValue := reflect.ValueOf(data)
	dataKind := dataValue.Type().Kind()

	if dataKind == reflect.Ptr {
		dataValue = dataValue.Elem()
		dataKind = dataValue.Type().Kind()
	}

	if dataKind != reflect.Slice && dataKind != reflect.Array {
		return nil, errors.New("value is not a slice or a array")
	}

	dataLen := dataValue.Len()
	removeCount := 0
	result := reflect.MakeSlice(reflect.SliceOf(dataValue.Type().Elem()), 0, 0)
	result = reflect.AppendSlice(result, dataValue)
	for i := 0; i < dataLen; i++ {
		if !condition(i, dataValue.Index(i).Interface()) {
			result = reflect.AppendSlice(result.Slice(0, i-removeCount), result.Slice(i+1-removeCount, result.Len()))
			removeCount++
		}
	}
	return result.Interface(), nil
}

func Any(condition func(interface{}) bool, args ...interface{}) bool {
	if len(args) == 1 {
		value := reflect.ValueOf(args[0])
		if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
			return Any(condition, prtArrayConvert(args[0])...)
		}
		if value.Kind() == reflect.Ptr {
			kind := value.Elem().Kind()
			if kind == reflect.Array || kind == reflect.Slice {
				return Any(condition, prtArrayConvert(args[0])...)
			}
		}
	}

	for _, item := range args {
		result := condition(item)
		if result {
			return result
		}
	}
	return false
}

func Each(action func(interface{}), args ...interface{}) {
	if len(args) == 1 {
		value := reflect.ValueOf(args[0])
		if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
			Each(action, prtArrayConvert(args[0])...)
			return
		}
		if value.Kind() == reflect.Ptr {
			kind := value.Elem().Kind()
			if kind == reflect.Array || kind == reflect.Slice {
				Each(action, prtArrayConvert(args[0])...)
				return
			}
		}
	}

	for _, item := range args {
		action(item)
	}
}

//*[]class/[]class -> []*interface{}
func prtArrayConvert(array interface{}) []interface{} {
	arrayValue := reflect.ValueOf(array)
	if arrayValue.Kind() == reflect.Ptr {
		arrayValue = arrayValue.Elem()
	}
	arrayLen := arrayValue.Len()
	var result = make([]interface{}, arrayLen)

	for i := 0; i < arrayLen; i++ {
		item := arrayValue.Index(i)

		if item.Kind() == reflect.Ptr {
			result[i] = item
		} else {
			result[i] = item.Addr().Interface()
		}
	}
	return result
}

// []class -> []interface{}
func ToArray(array interface{}) []interface{} {
	arrayValue := reflect.ValueOf(array)
	kind := arrayValue.Kind()

	if kind == reflect.Ptr {
		arrayValue = arrayValue.Elem()
	}

	arrayLen := arrayValue.Len()
	var result = make([]interface{}, arrayLen)

	for i := 0; i < arrayLen; i++ {
		item := arrayValue.Index(i).Interface()
		result[i] = item
	}
	return result
}
