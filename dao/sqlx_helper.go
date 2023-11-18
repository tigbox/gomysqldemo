package dao

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func getFields(obj interface{}) []string {
	res := make([]string, 0)
	v := reflect.TypeOf(obj)
	reflectValue := reflect.ValueOf(obj)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		field := reflectValue.Field(i).Interface()
		if v.Field(i).Type.Kind() == reflect.Struct && v.Field(i).Anonymous {
			res = append(res, getFields(field)...)
		} else {
			// skip empty field
			// https://forum.golangbridge.org/t/how-to-find-the-empty-field-in-struct-using-reflect/5819
			if !isEmptyValue(reflectValue) {
				tag := v.Field(i).Tag.Get("db")
				// skip auto incr field
				if tag == "id" {
					continue
				}
				res = append(res, tag)
			}
		}
	}
	return res
}

func getMap(obj interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("db"); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out
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

func genInsertSQL(c context.Context, table string, model interface{}) string {
	var dbFields, placeHolders []string
	dbFields = getFields(model)
	for _, fieldName := range dbFields {
		placeHolder := fmt.Sprintf(":%s", fieldName)
		placeHolders = append(placeHolders, placeHolder)
	}
	dbFieldStr := strings.Join(dbFields, ", ")
	placeHolderStr := strings.Join(placeHolders, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, dbFieldStr, placeHolderStr)
}

func (dao *MysqlDao) updateByCond(c context.Context, table string, model interface{}, cond interface{}) {
	clause := getMap(model)
	if len(clause) == 0 {
		return
	}
	sqlStr, args, err := sq.Update(table).
		SetMap(clause).
		Where(cond).ToSql()
	maybePanic(err)
	_, err = dao.db.Exec(sqlStr, args...)
	maybePanic(err)
}

func (dao *MysqlDao) selectBy(c context.Context, table string,
	selectExp, joinExp, orderBy []string, cond interface{},
	groupBy []string, offset, limit *int64, res interface{}) {
	builder := sq.Select(selectExp...).
		From(table)
	for _, exp := range joinExp {
		// TODO: support different kinds of join
		// left join
		// https://stackoverflow.com/questions/9770366/difference-in-mysql-join-vs-left-join
		builder = builder.LeftJoin(exp)
	}
	if cond != nil {
		builder = builder.Where(cond)
	}
	if groupBy != nil {
		builder = builder.GroupBy(groupBy...)
	}
	if orderBy != nil {
		builder = builder.OrderBy(orderBy...)
	}
	if offset != nil {
		builder = builder.Offset(uint64(*offset))
	}
	if limit != nil {
		builder = builder.Limit(uint64(*limit))
	}
	sqlStr, args, err := builder.ToSql()
	//spew.Dump(sqlStr)
	//spew.Dump(args)
	// fmt.Println(sqlStr)
	maybePanic(err)
	err = dao.db.Select(res, sqlStr, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		maybePanic(err)
	}
}

func (dao *MysqlDao) simpleSelect(c context.Context, table string, selectExp []string, cond, res interface{}) {
	dao.selectByCond(c, table, selectExp, nil, nil, cond, res)
}

// selectByCond get results by cond. query results are put in `res`
func (dao *MysqlDao) selectByCond(c context.Context, table string,
	selectExp, joinExp, orderBy []string, cond, res interface{}) {

	dao.selectBy(c, table, selectExp, joinExp, orderBy, cond, nil, nil, nil, res)
}
