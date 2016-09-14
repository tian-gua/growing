package gorm

import (
	"database/sql"
	"reflect"
	"github.com/aidonggua/growing/gutils"
	"fmt"
	"time"
)

var (
	//数据库连接对象
	gdb *sql.DB = nil
)

//插入或者更新一条记录
//插入和更新取决于 id 字段是否为0
func Save(obj interface{}, gtx ...*Transaction) error {
	var stmt *sql.Stmt
	var err error
	//生成sql
	sqlStr := parseSaveSql(obj)

	fmt.Println("[sql-gorm-" + gutils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

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
		return fmt.Errorf("no record changed");
	}
	return nil

}


//删除一条记录
func Delete(obj interface{}, gtx ...*Transaction) error {
	//生成sql
	sqlStr := ParseDeleteByPrimaryKeySql(obj)

	fmt.Println("[sql-gorm-" + gutils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

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
		return fmt.Errorf("no record changed");
	}
	return nil
}


//查询记录
func Query(param, resultSet interface{}, gtx ...*Transaction) error {
	pramValue := reflect.Indirect(reflect.ValueOf(param))
	rsValue := reflect.Indirect(reflect.ValueOf(resultSet))
	//make一个slice,因为可能resultSlice可能已经存在值
	//这里其实用不用都行,这里make一个或者用户自己new一个都OK
	//newSlice := reflect.MakeSlice(rsValue.Type(), 0, 0)
	sqlStr := ParseQuerySql(param)

	fmt.Println("[sql-gorm-" + gutils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

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
	//遍历所有记录
	for rows.Next() {
		err := rows.Scan(scans...)
		if err != nil {
			return err
		}
		//根据反射来新建一个和记录对应的对象
		var newV = reflect.New(pramValue.Type()).Elem()
		for i := 0; i < colNum; i++ {
			colName := columns[i]
			setRawData(newV.FieldByName(gutils.ToCamelCase(colName)), values[i])
		}
		rsValue = reflect.Append(rsValue, newV)
	}

	//更新target的值
	reflect.Indirect(reflect.ValueOf(resultSet)).Set(rsValue)
	return nil

}

//查询所有记录
func QueryAll(resultSet interface{}, gtx ...*Transaction) error {
	//获得target的反射信息
	resultsetRawData := reflect.Indirect(reflect.ValueOf(resultSet))
	//获得 切片元素 的反射信息
	element := getEmptySliceValue(resultsetRawData)
	elementType := element.Type()
	//生成sql
	sqlStr := ParseQueryAllSql(element.Interface())

	fmt.Println("[sql-gorm-" + gutils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

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
			setRawData(newV.FieldByName(gutils.ToCamelCase(colName)), values[i])
		}
		resultsetRawData = reflect.Append(resultsetRawData, newV)
	}
	//更新target的值
	reflect.ValueOf(resultSet).Elem().Set(resultsetRawData)
	return nil

}


//执行之定义sql查询语句
func CustomQuery(sqlStr string, resultSet interface{}, gtx ...*Transaction) error {
	fmt.Println("[sql-gorm-" + gutils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)
	//获得target的反射信息
	resultsetRawData := reflect.Indirect(reflect.ValueOf(resultSet))
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
	if reflect.Slice == resultsetRawData.Type().Kind() {

		elementType := getEmptySliceValue(resultsetRawData).Type()
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
					setRawData(newV.FieldByName(gutils.ToCamelCase(colName)), values[i])
				}
			} else {
				//如果是 基础类型 则直接赋值
				setRawData(newV, values[0])
			}
			resultsetRawData = reflect.Append(resultsetRawData, newV)
			index++
		}
		//更新target的值
		reflect.ValueOf(resultSet).Elem().Set(resultsetRawData)

	} else {
		//target为单条记录
		rows.Next()
		err := rows.Scan(scans...)
		if err != nil {
			return err
		}
		//如果target的为单个结构体
		if reflect.Struct == resultsetRawData.Type().Kind() {
			for i := 0; i < colNum; i++ {
				colName := columns[i]
				setRawData(resultsetRawData.FieldByName(gutils.ToCamelCase(colName)), values[i])
			}
		} else {
			setRawData(resultsetRawData, values[0])
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
