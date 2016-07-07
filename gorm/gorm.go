package gorm

import (
	"database/sql"
	"errors"
	"reflect"
	"github.com/aidonggua/growing/gutils"
)

//数据库连接对象
var gdb *sql.DB = nil


//插入或者更新一条记录
//插入和更新取决于 id 字段是否为0
func Save(obj interface{}, gtx ...*Transaction) error {
	var stmt *sql.Stmt
	var err error
	//生成sql
	sqlStr := parseSaveSql(obj)
	//判断是否在事务中执行
	if len(gtx) > 0 {
		stmt, err = gtx[0].tx.Prepare(sqlStr)
		if err != nil {
			return err
		}
	} else {
		//从sql.DB里获得stmt
		stmt, err = gdb.Prepare(sqlStr)
		if err != nil {
			return err
		}
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	if err != nil {
		return err
	}
	rownum, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rownum == 0 {
		return errors.New("no record insert")
	}

	return nil

}


//删除一条记录
func Delete(obj interface{}, gtx ...*Transaction) error {
	//生成sql
	sqlStr := parseSaveSql(obj)
	stmt, err := getStatement(sqlStr, gtx...)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	if err != nil {
		return err
	}
	rownum, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rownum == 0 {
		return errors.New("no record delete")
	}
	return nil
}


//查询记录
func Query(obj, target interface{}, gtx ...*Transaction) error {
	tv := reflect.Indirect(reflect.ValueOf(obj))
	t := tv.Type()
	targetVlaue := reflect.Indirect(reflect.ValueOf(target))
	sqlStr := parseQuerySql(obj)
	stmt, err := getStatement(sqlStr, gtx...)
	if err != nil {
		return err
	}
	defer stmt.Close()
	//查询
	rows, err := stmt.Query()
	if err != nil {
		return err
	}
	defer rows.Close()
	//获得所有列
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	//获得列的数量
	colNum := len(columns)

	values := make([]sql.RawBytes, colNum)
	scans := make([]interface{}, colNum)
	//封装
	for i := range values {
		scans[i] = &values[i]
	}
	var index int = 0
	//遍历所有记录
	for rows.Next() {
		err := rows.Scan(scans...)
		if err != nil {
			return err
		}
		//根据反射来新建一个和记录对应的对象
		var newV = reflect.New(t).Elem()
		for i := 0; i < colNum; i++ {
			colName := columns[i]
			setValue(newV.FieldByName(gutils.ToCamelCase(colName)), values[i])
		}
		targetVlaue = reflect.Append(targetVlaue, newV)
		index++
	}
	//更新target的值
	reflect.ValueOf(target).Elem().Set(targetVlaue.Slice(0, index))
	return nil

}

//查询所有记录
func QueryAll(target interface{}, gtx ...*Transaction) error {
	//获得target的反射信息
	targetV := reflect.Indirect(reflect.ValueOf(target))
	//获得 切片元素 的反射信息
	element := getEmptySliceValue(targetV)
	elementType := element.Type()
	//生成sql
	sqlStr := parseQueryAllSql(element.Interface())
	stmt, err := getStatement(sqlStr, gtx...)
	if err != nil {
		return err
	}
	defer stmt.Close()
	//查询
	rows, err := stmt.Query()
	if err != nil {
		return err
	}
	defer rows.Close()
	//获得所有列
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	//获得列的数量
	colNum := len(columns)

	values := make([]sql.RawBytes, colNum)
	scans := make([]interface{}, colNum)
	//封装
	for i := range values {
		scans[i] = &values[i]
	}
	var index int = 0
	//遍历所有记录
	for rows.Next() {
		err := rows.Scan(scans...)
		if err != nil {
			return err
		}
		//根据反射来新建一个和记录对应的对象
		var newV = reflect.New(elementType).Elem()
		for i := 0; i < colNum; i++ {
			colName := columns[i]
			setValue(newV.FieldByName(gutils.ToCamelCase(colName)), values[i])
		}
		targetV = reflect.Append(targetV, newV)
		index++
	}
	//更新target的值
	reflect.ValueOf(target).Elem().Set(targetV.Slice(0, index))
	return nil

}


//执行之定义sql查询语句
func CustomQuery(sqlStr string, target interface{}, gtx ...*Transaction) error {
	//获得target的反射信息
	targetV := reflect.Indirect(reflect.ValueOf(target))
	stmt, err := getStatement(sqlStr, gtx...)
	if err != nil {
		return err
	}
	defer stmt.Close()
	//查询
	rows, err := stmt.Query()
	if err != nil {
		return err
	}
	defer rows.Close()
	//获得所有列
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	//获得列的数量
	colNum := len(columns)

	values := make([]sql.RawBytes, colNum)
	scans := make([]interface{}, colNum)
	//封装
	for i := range values {
		scans[i] = &values[i]
	}
	//如果传过来的是一个切片
	if reflect.Slice == targetV.Type().Kind() {
		elementType := getEmptySliceValue(targetV).Type()
		var index int = 0
		//遍历所有记录
		for rows.Next() {
			err := rows.Scan(scans...)
			if err != nil {
				return err
			}

			//根据反射来新建一个和记录对应的对象
			var newV = reflect.New(elementType).Elem()
			//如果切片类型为结构体
			if reflect.Struct == elementType.Kind() {
				for i := 0; i < colNum; i++ {
					colName := columns[i]
					setValue(newV.FieldByName(gutils.ToCamelCase(colName)), values[i])
				}
			} else {
				//如果是 基础类型 则直接赋值
				setValue(newV, values[0])
			}
			targetV = reflect.Append(targetV, newV)
			index++
		}
		//更新target的值
		reflect.ValueOf(target).Elem().Set(targetV.Slice(0, index))

	} else {
		//target为单条记录
		rows.Next()
		err := rows.Scan(scans...)
		if err != nil {
			return err
		}
		//如果target的为单个结构体
		if reflect.Struct == targetV.Type().Kind() {
			for i := 0; i < colNum; i++ {
				colName := columns[i]
				setValue(targetV.FieldByName(gutils.ToCamelCase(colName)), values[i])
			}
		} else {
			setValue(targetV, values[0])
		}
	}
	return nil

}

//获得空切片元素的类型
func getEmptySliceValue(slice reflect.Value) reflect.Value {
	t := slice.Type()
	//给切片元素开辟一个空间
	vSlice := reflect.MakeSlice(t, 1, 1)
	//获得 切片元素 的反射信息
	element := vSlice.Slice(0, 1).Index(0)
	return element
}

//获得statement,有事务和非事务2种情况
func getStatement(sqlStr string, gtx ...*Transaction) (*sql.Stmt, error) {
	var stmt *sql.Stmt
	var err error
	//判断是否在事务中执行
	if len(gtx) > 0 {
		stmt, err = gtx[0].tx.Prepare(sqlStr)
		if err != nil {
			return nil, err
		}
	} else {
		//从sql.DB里获得stmt
		stmt, err = gdb.Prepare(sqlStr)
		if err != nil {
			return nil, err
		}
	}
	return stmt, err
}


//关闭DB对象
func CloseDB() {
	gdb.Close()
}
