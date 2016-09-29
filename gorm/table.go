package gorm

//Table接口用于获得结构体的指定table名
//例如 type abc struct 如果不实现Table接口
//则默认table名为abc，但是实际上abc映射的表可能为tb_abc
//GetTableName()函数可以设置指定表名
type ITable interface {
	GetTableName() string
}
