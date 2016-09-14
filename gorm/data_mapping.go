package gorm

import (
	"reflect"
	"fmt"
	"strings"
)

var StructInfoMap = make(map[reflect.Type]*StructInfo)

//结构体信息
type StructInfo struct {
	FieldsMap map[string]*StructField //字段字典集合
	Name      string                  //类型名
	TableName string                  //表名
}

//结构体字段信息
type StructField struct {
	name           string        //字段名
	value          reflect.Value //字段值
	stringValue    string        //字符串值
	tableFieldName string        //表属性名
	tableFieldType string        //表属性类型
}


//获得结构体的信息
func GetStructInfo(target interface{}) *StructInfo {
	v := reflect.Indirect(reflect.ValueOf(target))
	t := v.Type()
	//判断target的类型
	if t.Kind() != reflect.Struct {
		fmt.Println("not struct")
		return nil
	}
	return GetReflectInfo(t, v)
}

//获得结构体的反射的信息
func GetReflectInfo(t reflect.Type, v reflect.Value) *StructInfo {

	var structInfo *StructInfo

	fieldsMap := make(map[string]*StructField)
	//从map里取结构体信息,如果map没有则新建一个然后存map
	if value, ok := StructInfoMap[t]; ok {
		structInfo = value
		//更新缓存的结构体字段的值,这一部分肯定不能使用缓存,因为sql的条件都不同
		for key, _ := range structInfo.FieldsMap {
			//更新字段的value属性
			structInfo.FieldsMap[key].value.Set(v.FieldByName(key))
			//更新字段的stringValue属性
			structInfo.FieldsMap[key].stringValue = parseValueToDBString(v.FieldByName(key))
		}
	} else {
		//遍历所有属性
		for index := 0; index < t.NumField(); index++ {
			structField := t.Field(index)
			structFieldValue := v.Field(index)

			//获取field标签的值 作为数据库字段名
			tableField := strings.TrimSpace(structField.Tag.Get("field"))

			//如果字段
			if len(tableField) != 0 {
				//构造一个新的StructField
				sf := &StructField{
					name:structField.Name,
					tableFieldName:tableField,
					tableFieldType:getDataType(t.Kind().String()),
					value:structFieldValue,
					stringValue:parseValueToDBString(structFieldValue),
				}
				//将新的StructField放入Map
				fieldsMap[tableField] = sf
			}
		}
		//构造一个新的StructInfo
		structInfo = &StructInfo{
			Name:t.Name(),
			TableName:getTableName(t.Name()),
			FieldsMap:fieldsMap,
		}
		//将新的StructInfo放入Map当缓存用
		StructInfoMap[t] = structInfo
	}
	return structInfo
}


