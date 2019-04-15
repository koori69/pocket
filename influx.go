// Package core @author K·J Create at 2019-04-15 15:04
package core

import (
	"encoding/json"
	"errors"
	"reflect"
	"unsafe"
)

// Stom convert []string to map[int]string
func Stom(columns []string) map[int]string {
	if nil == columns || len(columns) == 0 {
		return nil
	}
	index := make(map[int]string, 0)
	for i, v := range columns {
		index[i] = v
	}
	return index
}

// UnmarshalInflux convert influx db row to struct
// columns [index]:[column name]
// data row of values
// v struct tag:'influx'
func UnmarshalInflux(columns map[int]string, data []interface{}, v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return errors.New("interface must be a ptr")
	}
	value := reflect.ValueOf(v).Elem()
	for i, k := range data {
		tag := columns[i]
		for j := 0; j < t.Elem().NumField(); j++ {
			if tag == t.Elem().Field(i).Tag.Get("influx") {
				err := setValue(value.Field(i), k)
				if nil != err {
					return err
				}
			}
		}
	}
	return nil
}

// setValue set value
func setValue(v reflect.Value, data interface{}) error {
	if nil == data {
		return errors.New("param invalid")
	}
	if v.CanSet() {
		switch v.Kind() {
		case reflect.Ptr:
			switch v.Type().String() {
			case "*int":
				val, err := (data.(json.Number)).Int64()
				if nil != err {
					return err
				}
				p := new(int)
				*p = int(val)
				v.Set(reflect.ValueOf(p))
			case "*int64":
				val, err := (data.(json.Number)).Int64()
				if nil != err {
					return err
				}
				p := new(int64)
				*p = val
				v.Set(reflect.ValueOf(p))
			case "*string":
				p := new(string)
				*p = data.(string)
				v.Set(reflect.ValueOf(p))
			default:
				return errors.New("not support type")
			}
		case reflect.String:
			v.SetString(data.(string))
		case reflect.Int:
			fallthrough
		case reflect.Int64:
			val, err := (data.(json.Number)).Int64()
			if nil != err {
				return err
			}
			v.SetInt(val)
		default:
			return errors.New("not support type")
		}
	} else {
		if !v.CanAddr() {
			return errors.New("can not get addr")
		}
		addr := v.Addr()
		ptr := unsafe.Pointer(addr.Pointer())
		switch v.Type().String() {
		case "int":
			fallthrough
		case "*int":
			val, err := (data.(json.Number)).Int64()
			if nil != err {
				return err
			}
			*(*int)(ptr) = int(val)
		case "int64":
			fallthrough
		case "*int64":
			val, err := (data.(json.Number)).Int64()
			if nil != err {
				return err
			}
			*(*int64)(ptr) = val
		case "*string":
			fallthrough
		case "string":
			*(*string)(ptr) = data.(string)
		default:
			return errors.New("not support type")
		}
	}
	return nil
}
