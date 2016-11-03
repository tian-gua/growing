package gorm

import (
	"database/sql"
	"fmt"
	"github.com/aidonggua/growing/gutils"
	"reflect"
	"strings"
	"time"
)

//插入或者更新一条记录
//插入和更新取决于 id 字段是否为0
func Save(obj interface{}, gtx ...*Transaction) (int64, error) {
	var err error
	//生成sql
	sqlStr, err := ParseSaveSql(obj)
	if err != nil {
		return 0, err
	}
	fmt.Println("[sql-gorm-" + gutils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

	stmt, err := getStatement(sqlStr, gtx...)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		return 0, err
	}
	if strings.HasPrefix(sqlStr, "insert") {
		return result.LastInsertId()
	}
	return result.RowsAffected()
}

/**有选择的插入记录，如果字段为零值则不插入*/
func InsertSelective(obj interface{}, gtx ...*Transaction) (int64, error) {
	return insertOrupdate(obj, true, true, gtx...)
}

/**插入记录*/
func Insert(obj interface{}, gtx ...*Transaction) (int64, error) {
	return insertOrupdate(obj, false, true, gtx...)
}

/**有选择的更新记录，如果字段为零值则不更新*/
func UpdateSelective(obj interface{}, gtx ...*Transaction) (int64, error) {
	return insertOrupdate(obj, true, false, gtx...)
}

/**更新记录*/
func Update(obj interface{}, gtx ...*Transaction) (int64, error) {
	return insertOrupdate(obj, false, false, gtx...)
}

func insertOrupdate(obj interface{}, isSelective bool, isInsert bool, gtx ...*Transaction) (int64, error) {
	var err error
	var sqlStr string

	/**生成sql语句*/
	if isInsert {
		sqlStr, err = ParseInsertSql(obj, isSelective)
	} else {
		sqlStr, err = ParseUpdateSql(obj, isSelective)
	}
	if err != nil {
		return 0, err
	}
	fmt.Println("[sql-gorm-" + gutils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

	/**获得声明*/
	stmt, err := getStatement(sqlStr, gtx...)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	/**执行sql*/
	result, err := stmt.Exec()
	if err != nil {
		return 0, err
	}

	/**如果是insert 返回插入记录的id 否则返回 更新的条数*/
	if isInsert {
		return result.LastInsertId()
	}
	return result.RowsAffected()
}

//删除一条记录
func Delete(obj interface{}, gtx ...*Transaction) (int64, error) {
	//生成sql
	sqlStr, err := ParseDeleteByPrimaryKeySql(obj)
	if err != nil {
		return 0, err
	}
	fmt.Println("[sql-gorm-" + gutils.DateFormat(time.Now(), "yyyy-MM-dd HH:mm:ss") + "]:" + sqlStr)

	stmt, err := getStatement(sqlStr, gtx...)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//查询记录
func Query(param, resultSet interface{}, gtx ...*Transaction) error {
	pramValue := reflect.Indirect(reflect.ValueOf(param))
	resultSetData := reflect.Indirect(reflect.ValueOf(resultSet))

	sqlStr, err := ParseQuerySql(param)
	if err != nil {
		return err
	}
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
		//根据反射来新建一个临时value和记录对应的对象
		var temporaryValue = reflect.New(pramValue.Type()).Elem()
		for i := 0; i < colNum; i++ {
			colName := columns[i]
			setRawData(temporaryValue.FieldByName(toCamelCase(colName)), values[i])
		}
		resultSetData = reflect.Append(resultSetData, temporaryValue)
	}

	//更新target的值
	reflect.Indirect(reflect.ValueOf(resultSet)).Set(resultSetData)
	return nil
}

//查询所有记录
func QueryAll(resultSet interface{}, gtx ...*Transaction) error {
	//获得target的反射信息
	resultsetRawData := reflect.Indirect(reflect.ValueOf(resultSet))
	if resultsetRawData.Type().Kind() != reflect.Slice {
		return fmt.Errorf("not slice param")
	}
	//获取切片的元素的类型
	elementType := resultsetRawData.Type().Elem()
	//生成sql
	sqlStr, err := ParseQueryAllSql(reflect.New(elementType).Interface())
	if err != nil {
		return err
	}
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
		//根据反射来新建一个临时value和记录对应的对象
		var temporaryValue = reflect.New(elementType).Elem()
		for i := 0; i < colNum; i++ {
			colName := columns[i]
			setRawData(temporaryValue.FieldByName(toCamelCase(colName)), values[i])
		}
		resultsetRawData = reflect.Append(resultsetRawData, temporaryValue)
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
		//获取切片
		elementType := resultsetRawData.Type().Elem()
		var index int = 0
		//遍历所有记录
		for rows.Next() {
			err := rows.Scan(scans...)
			if err != nil {
				return err
			}
			//根据反射来新建一个临时value和记录对应的对象
			var temporaryValue = reflect.New(elementType).Elem()
			//如果切片类型为结构体
			if reflect.Struct == elementType.Kind() {
				for i := 0; i < colNum; i++ {
					colName := columns[i]
					temporaryV := temporaryValue.FieldByName(toCamelCase(colName))
					//判断结构体是否有对应的属性，或者该属性是否可修改
					//如果判断失败则跳过此属性的复制
					if temporaryV.CanSet() {
						setRawData(temporaryV, values[i])
					}
				}
			} else {
				//如果是 基础类型 则直接赋值
				setRawData(temporaryValue, values[0])
			}
			resultsetRawData = reflect.Append(resultsetRawData, temporaryValue)
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

//获得statement,有事务和非事务2种情况
func getStatement(sqlStr string, gtx ...*Transaction) (*sql.Stmt, error) {
	//校验是否初始化
	if !isInit {
		panic("no db init")
	}

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
