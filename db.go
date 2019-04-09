// Package pocket @author K·J Create at 2019-04-09 15:14
package pocket

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// GenerateUpdateSet get set query
func GenerateUpdateSet(origin interface{}) (string, []interface{}) {
	set := ""
	params := make([]interface{}, 0)
	originType := reflect.TypeOf(origin)
	if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
		log.Println("param error")
		return set, params
	}
	originValue := reflect.ValueOf(origin)

	for i := 0; i < originType.Elem().NumField(); i++ {
		tag := originType.Elem().Field(i).Tag.Get("db_opt")
		if strings.Contains(tag, "set") {
			if !originValue.Elem().Field(i).IsNil() {
				if "" != set {
					set += ","
				}
				set += "`" + originType.Elem().Field(i).Tag.Get("db_column") + "`=?"
				params = append(params, originValue.Elem().Field(i).Interface())
			}
		}
	}
	return set, params
}

// GenerateAdd add record
func GenerateAdd(origin interface{}) (string, []interface{}) {
	column := ""
	values := ""
	params := make([]interface{}, 0)
	originType := reflect.TypeOf(origin)
	if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
		log.Println("param error")
		return "", params
	}
	originValue := reflect.ValueOf(origin)

	for i := 0; i < originType.Elem().NumField(); i++ {
		tag := originType.Elem().Field(i).Tag.Get("db_opt")
		if strings.Contains(tag, "add") {
			if !originValue.Elem().Field(i).IsNil() {
				if "" != column {
					column += ","
					values += ","
				}
				column += "`" + originType.Elem().Field(i).Tag.Get("db_column") + "`"
				values += "?"
				params = append(params, originValue.Elem().Field(i).Interface())
			}
		}
	}
	if len(column) == 0 {
		return "", params
	}
	return fmt.Sprintf("(%s) VALUES (%s)", column, values), params
}

// GenerateValues gen values
func GenerateValues(origin interface{}) (string, []interface{}) {
	values := ""
	params := make([]interface{}, 0)
	originType := reflect.TypeOf(origin)
	if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
		log.Println("param error")
		return "", params
	}
	originValue := reflect.ValueOf(origin)

	for i := 0; i < originType.Elem().NumField(); i++ {
		tag := originType.Elem().Field(i).Tag.Get("db_opt")
		if strings.Contains(tag, "add") {
			if !originValue.Elem().Field(i).IsNil() {
				if "" != values {
					values += ","
				}
				values += "?"
				params = append(params, originValue.Elem().Field(i).Interface())
			}
		}
	}
	if len(values) == 0 {
		return "", params
	}
	return fmt.Sprintf("(%s)", values), params
}

// GenerateSort order str
func GenerateSort(sort string, origin interface{}) string {
	order := ""
	originType := reflect.TypeOf(origin)
	if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct || "" == sort {
		log.Println("param error")
		return ""
	}

	list := strings.Split(sort, ",")
	for _, v := range list {
		desc := false
		if strings.HasPrefix(v, "-") {
			desc = true
			v = v[1:]
		}
		for i := 0; i < originType.Elem().NumField(); i++ {
			obj := originType.Elem().Field(i)
			if obj.Tag.Get("db_column") == v {
				tag := obj.Tag.Get("db_opt")
				if strings.Contains(tag, "sort") {
					if desc {
						order += obj.Tag.Get("db_column") + " DESC,"
					} else {
						order += obj.Tag.Get("db_column") + ","
					}
					break
				}
			}
		}
	}
	if len(order) > 0 {
		order = order[0 : len(order)-1]
	}
	return order
}
