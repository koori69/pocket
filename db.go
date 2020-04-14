// Package pocket @author K·J Create at 2019-04-09 15:14
package pocket

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	generateSqlErr = "生成sql语句错误"
	columnErr      = "列数不一致"
)

type SqlBuilder struct {
	Table string
	Model interface{}
}

// NewSqlBuilder sql builder tag规则 `db:"column,add,set,sort"`
func NewSqlBuilder(table string, model interface{}) *SqlBuilder {
	return &SqlBuilder{table, model}
}

// BuildInsertRow 生成单条插入sql
func (builder *SqlBuilder) BuildInsertRow() (string, []interface{}, error) {
	sql, param := builder.generate("add-row")
	if "" == sql {
		DefaultLogger.Error(generateSqlErr)
		return "", nil, errors.New(generateSqlErr)
	}

	return "INSERT INTO `" + builder.Table + "`" + sql, param, nil
}

// BuildInsert 生成批量插入sql
func (builder *SqlBuilder) BuildInsert() (string, []interface{}, error) {
	sql, param := builder.generate("add-rows")
	if "" == sql {
		DefaultLogger.Error(generateSqlErr)
		return "", nil, errors.New(generateSqlErr)
	}

	return "INSERT INTO " + builder.Table + sql, param, nil
}

// generate 生成sql的列及条件部分，返回 列，条件部分
func (builder *SqlBuilder) generate(action string) (string, []interface{}) {
	switch action {
	case "add-row":
		column := ""
		values := ""
		params := make([]interface{}, 0)
		originType := reflect.TypeOf(builder.Model)
		if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
			DefaultLogger.Warn("param error")
			return "", params
		}
		originValue := reflect.ValueOf(builder.Model)

		for i := 0; i < originType.Elem().NumField(); i++ {
			tag := originType.Elem().Field(i).Tag.Get("db")
			if "" == tag {
				continue
			}
			if originValue.Elem().Field(i).Kind() == reflect.Ptr {
				if !originValue.Elem().Field(i).IsNil() {
					if "" != column {
						column += ","
						values += ","
					}
					if strings.Index(tag, ",") > 0 {
						column += "`" + tag[:strings.Index(tag, ",")] + "`"
					} else {
						column += "`" + tag + "`"
					}
					values += "?"
					params = append(params, originValue.Elem().Field(i).Interface())
				}
			}
		}
		if len(column) == 0 {
			return "", params
		}
		return fmt.Sprintf("(%s) VALUES (%s)", column, values), params
	case "add-rows":
		column := ""
		columnCount := make(map[int]int, 10)
		values := ""
		params := make([]interface{}, 0)
		originType := reflect.TypeOf(builder.Model)
		if originType.Kind() != reflect.Slice && originType.Kind() != reflect.Array ||
			originType.Elem().Kind() != reflect.Struct {
			DefaultLogger.Warn("param error")
			return "", params
		}
		originValue := reflect.ValueOf(builder.Model)
		for j := 0; j < originValue.Len(); j++ {
			item := originValue.Index(j)
			itemType := item.Type()
			row := ""
			for i := 0; i < itemType.NumField(); i++ {
				tag := itemType.Field(i).Tag.Get("db")
				if "" == tag {
					continue
				}
				if item.Field(i).Kind() == reflect.Ptr {
					if !item.Field(i).IsNil() {
						if j == 0 {
							if i > 0 {
								column += ","
							}
							if strings.Index(tag, ",") > 0 {
								column += "`" + tag[:strings.Index(tag, ",")] + "`"
							} else {
								column += "`" + tag + "`"
							}
						}

						if num, ok := columnCount[i]; ok {
							columnCount[i] = num + 1
						} else {
							columnCount[i] = 1
						}
						if "" != row {
							row += ","
						}

						row += "?"
						params = append(params, item.Field(i).Interface())
					}
				}
			}
			if len(row) > 0 {
				if j > 0 {
					values += ","
				}
				values += fmt.Sprintf("(%s)", row)
			}
		}

		if len(values) == 0 {
			DefaultLogger.Warn("no data")
			return "", params
		}
		index := 0
		count := 0
		for _, v := range columnCount {
			if index == 0 {
				count = v
			}
			if count != v {
				DefaultLogger.Error(columnErr)
				return "", nil
			}
			index++
		}
		return fmt.Sprintf("(%s) VALUES %s", column, values), params
	case "set":
		set := ""
		params := make([]interface{}, 0)
		originType := reflect.TypeOf(builder.Model)
		if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
			DefaultLogger.Warn("param error")
			return set, params
		}
		originValue := reflect.ValueOf(builder.Model)

		for i := 0; i < originType.Elem().NumField(); i++ {
			tag := originType.Elem().Field(i).Tag.Get("db")
			if "" == tag {
				continue
			}
			if strings.Contains(tag, "set") {
				if !originValue.Elem().Field(i).IsNil() {
					if "" != set {
						set += ","
					}
					if strings.Index(tag, ",") > 0 {
						set += "`" + tag[:strings.Index(tag, ",")] + "`=?"
					} else {
						set += "`" + tag + "`=?"
					}
					params = append(params, originValue.Elem().Field(i).Interface())
				}
			}
		}
		return set, params
	default:
		return "", nil
	}
	return "", nil
}

func XormUpdateParam(model interface{}) (map[string]interface{}, error) {
	params := make(map[string]interface{}, 0)
	originType := reflect.TypeOf(model)
	if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
		DefaultLogger.Warn("param error")
		return nil, errors.New("param error")
	}
	originValue := reflect.ValueOf(model)

	for i := 0; i < originType.Elem().NumField(); i++ {
		if !originValue.Elem().Field(i).IsNil() {
			tag := originType.Elem().Field(i).Tag.Get("db")
			if "" == tag {
				continue
			}
			if strings.Index(tag, ",") > 0 {
				params[tag[:strings.Index(tag, ",")]] = originValue.Elem().Field(i).Interface()
			} else {
				params[tag] = originValue.Elem().Field(i).Interface()
			}
		}
	}
	if len(params) == 0 {
		return nil, nil
	}
	return params, nil
}

// ReqToSql 查询列表时，将req转换为where查询条件及sort排序，仅支持单表
func ReqToSql(model interface{}) (string, []interface{}, map[string][]string, error) {
	originType := reflect.TypeOf(model)
	if originType.Kind() != reflect.Struct {
		DefaultLogger.Warn("param error")
		return "", nil, nil, errors.New("param error")
	}
	originValue := reflect.ValueOf(model)

	sort := make(map[string][]string, 0)
	param := make([]interface{}, 0)
	where := ""
	for i := 0; i < originType.NumField(); i++ {
		t := originType.Field(i).Type.Name()
		if originValue.Field(i).Kind() == reflect.Ptr {
			t = originType.Field(i).Type.String()
		}
		sqlTag := originType.Field(i).Tag.Get("sql")
		if "" != sqlTag {
			continue
		}
		switch t {
		case "string":
			if "" == originValue.Field(i).String() {
				continue
			}
			// 排除特殊字段
			if "Token" == originType.Field(i).Name {
				continue
			}
			if "Sort" == originType.Field(i).Name {
				sorts := strings.Split(originValue.Field(i).String(), ",")
				for _, v := range sorts {
					if strings.HasPrefix(v, "-") {
						if desc, ok := sort["desc"]; ok {
							sort["desc"] = append(desc, "`"+v[1:]+"`")
						} else {
							sort["desc"] = []string{"`" + v[1:] + "`"}
						}
					} else {
						if asc, ok := sort["asc"]; ok {
							sort["asc"] = append(asc, "`"+v+"`")
						} else {
							sort["asc"] = []string{"`" + v + "`"}
						}
					}
				}
				continue
			}
			if len(where) > 0 {
				where += " AND "
			}
			where += fmt.Sprintf("`%s`=?", SnakeString(originType.Field(i).Name))
			param = append(param, originValue.Field(i).String())
		case "*string":
			if originValue.Field(i).IsNil() {
				continue
			}
			value := getValue(originValue.Field(i).Elem())
			if "" == value {
				continue
			}
			// 排除特殊字段
			if "Token" == originType.Field(i).Name {
				continue
			}

			if "Sort" == originType.Field(i).Name {
				sorts := strings.Split(value, ",")
				for _, v := range sorts {
					if strings.HasPrefix(v, "-") {
						if desc, ok := sort["desc"]; ok {
							sort["desc"] = append(desc, "`"+v[1:]+"`")
						} else {
							sort["desc"] = []string{"`" + v[1:] + "`"}
						}
					} else {
						if asc, ok := sort["asc"]; ok {
							sort["asc"] = append(asc, "`"+v+"`")
						} else {
							sort["asc"] = []string{"`" + v + "`"}
						}
					}
				}
				continue
			}
			if len(where) > 0 {
				where += " AND "
			}
			where += fmt.Sprintf("`%s`=?", SnakeString(originType.Field(i).Name))
			param = append(param, value)
		case "*int", "*int64":
			if originValue.Field(i).IsNil() {
				continue
			}
			if len(where) > 0 {
				where += " AND "
			}
			if "BeginAt" == originType.Field(i).Name {
				where += "`created_at`>=?"
				param = append(param, getValue(originValue.Field(i).Elem()))
				continue
			}
			if "EndAt" == originType.Field(i).Name {
				where += "`created_at`<=?"
				param = append(param, getValue(originValue.Field(i).Elem()))
				continue
			}
			where += fmt.Sprintf("`%s`=?", SnakeString(originType.Field(i).Name))
			param = append(param, getValue(originValue.Field(i).Elem()))
		case "int", "int64":
			// 排除特殊字段
			if "Page" == originType.Field(i).Name || "Rows" == originType.Field(i).Name {
				continue
			}
			if len(where) > 0 {
				where += " AND "
			}
			if "BeginAt" == originType.Field(i).Name {
				where += "`created_at`>=?"
				param = append(param, originValue.Field(i).Int())
				continue
			}
			if "EndAt" == originType.Field(i).Name {
				where += "`created_at`<=?"
				param = append(param, originValue.Field(i).Int())
				continue
			}
			where += fmt.Sprintf("`%s`=?", SnakeString(originType.Field(i).Name))
			param = append(param, originValue.Field(i).Int())
		default:
			DefaultLogger.Warn("not supported type")
			continue
		}
	}
	return where, param, sort, nil
}
