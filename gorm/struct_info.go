package gorm

import (
	"reflect"
	"github.com/aidonggua/growing/gutils"
	"fmt"
)

var StructInfoMap = make(map[reflect.Type]*StructInfo)

//结构体信息
type StructInfo struct {
	fields          []StructField
	name            string
	tableName       string
	tableFieldNames []string
	fieldNames      []string
	subStructInfo   []StructInfo
}

//结构体字段信息
type StructField struct {
	name           string
	tableFieldName string
	tableFieldType string
}

func (this *StructInfo) Get() {
	fmt.Printf("%#v", this)
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
	//从map里取结构体信息,如果map没有则新建一个然后存map
	if value, ok := StructInfoMap[t]; ok {
		info = value
	} else {
		tableName := gutils.UnCamelCase(t.Name())
		tableFieldNames := new([]string)
		subStructInfo := new([]StructInfo)
		fields := new([]StructField)
		fieldNames := new([]string)
		for index := 0; index < t.NumField(); index++ {
			subsf := t.Field(index)
			subt := subsf.Type
			if subt.Kind() == reflect.Struct {
				//递归获得子结构体信息
				*subStructInfo = append(*subStructInfo, *GetValueInfo(subt))
			} else {
				sf := StructField{
					name:subsf.Name,
					tableFieldName:gutils.UnCamelCase(subsf.Name),
					tableFieldType:gutils.GetDBType(subt.Kind().String())}
				*fields = append(*fields, sf)
				*fieldNames = append(*fieldNames, subsf.Name)
				*tableFieldNames = append(*tableFieldNames, gutils.UnCamelCase(subsf.Name))
			}

		}

		info = &StructInfo{name:t.Name(), tableName:tableName, fields:*fields, tableFieldNames:*tableFieldNames, fieldNames:*fieldNames, subStructInfo:*subStructInfo }
		StructInfoMap[v.Type()] = info
	}

	return info

}

//获得结构体的反射的信息
func GetValueInfo(t reflect.Type) *StructInfo {
	var info *StructInfo
	//从map里取结构体信息,如果map没有则新建一个然后存map
	if value, ok := StructInfoMap[t]; ok {
		info = value
	} else {
		tableName := gutils.UnCamelCase(t.Name())
		tableFieldNames := new([]string)
		subStructInfo := new([]StructInfo)
		fields := new([]StructField)
		fieldNames := new([]string)

		for index := 0; index < t.NumField(); index++ {
			sf := t.Field(index)
			t := sf.Type
			if t.Kind() == reflect.Struct {
				//递归获得子结构体信息
				*subStructInfo = append(*subStructInfo, *GetValueInfo(t))
			} else {
				sf := StructField{
					name:sf.Name,
					tableFieldName:gutils.UnCamelCase(sf.Name),
					tableFieldType:gutils.GetDBType(t.Kind().String())}
				*fields = append(*fields, sf)
				*fieldNames = append(*fieldNames, sf.name)
				*tableFieldNames = append(*tableFieldNames, gutils.UnCamelCase(sf.name))
			}

		}
		info = &StructInfo{name:t.Name(), tableName:tableName, fields:*fields, tableFieldNames:*tableFieldNames, fieldNames:*fieldNames, subStructInfo:*subStructInfo }
		StructInfoMap[t] = info
	}
	return info
}
