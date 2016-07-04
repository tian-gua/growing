package gorm

import (
	"reflect"
	"github.com/aidonggua/growing/gutils"
	"fmt"
)

var StructInfoMap = make(map[reflect.Type]*StructInfo)

//结构体信息
type StructInfo struct {
	FieldsMap       map[string]StructField //字段字典集合
	Name            string                 //类型名
	TableName       string                 //表名
	TableFieldNames []string               //表字段集合
	FieldNames      []string               //属性名称集合
	SubStructInfo   []StructInfo           //子结构体
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

	var info *StructInfo
	v := reflect.Indirect(reflect.ValueOf(target))
	t := v.Type()
	//判断target的类型
	if t.Kind() != reflect.Struct {
		fmt.Println(GET_STRUCTINFO_ERROR)
		return nil
	}

	info = GetReflectInfo(t, v)

	return info

}

//获得结构体的反射的信息
func GetReflectInfo(t reflect.Type, v reflect.Value) *StructInfo {

	var haveCache bool = false
	var cache *StructInfo
	var info *StructInfo
	var tableName string
	var tName string
	tableFieldNames := new([]string)
	fieldNames := new([]string)
	subStructInfo := new([]StructInfo)
	fieldsMap := make(map[string]StructField)

	//从map里取结构体信息,如果map没有则新建一个然后存map
	if value, ok := StructInfoMap[t]; ok {
		haveCache = true
		cache = value
	}

	//遍历所有属性
	for index := 0; index < t.NumField(); index++ {
		sf := t.Field(index)
		sfv := v.Field(index)
		t := sf.Type
		//判断属性是否为结构体
		if sf.Type.Kind() == reflect.Struct {
			//递归获得子结构体信息
			*subStructInfo = append(*subStructInfo, *GetReflectInfo(sf.Type, sfv))
		} else {

			//如果有缓存则只更新StructField的value
			if !haveCache {
				sf := StructField{
					name:sf.Name,
					tableFieldName:gutils.UnCamelCase(sf.Name),
					tableFieldType:gutils.GetDBType(t.Kind().String()),
					value:sfv,
					stringValue:gutils.ParseValueToString(sfv)}
				fieldsMap[sf.name] = sf

				*fieldNames = append(*fieldNames, sf.name)
				*tableFieldNames = append(*tableFieldNames, gutils.UnCamelCase(sf.name))
			} else {
				fieldsMap[sf.Name].value = sfv
			}
		}
	}

	if haveCache {
		tableName = cache.TableName
		*fieldNames = cache.FieldNames
		*tableFieldNames = cache.TableFieldNames
		tName = cache.Name
	} else {
		tableName = gutils.UnCamelCase(t.Name())
		tName = t.Name()
	}

	info = &StructInfo{Name:tName, TableName:tableName, FieldsMap:fieldsMap, TableFieldNames:*tableFieldNames, FieldNames:*fieldNames, SubStructInfo:*subStructInfo }
	StructInfoMap[t] = info

	return info
}
