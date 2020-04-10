package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

// 使用 dialect 隔离不同数据库之间的差异，便于扩展。
// Dialect is an interface contains methods that a dialect has to implement
type Dialect interface {
	DataTypeOf(typ reflect.Value) string  			       // 用于将 Go 语言的类型转换为该数据库的数据类型。
	TableExistSQL(tableName string) (string, []interface{})// 返回某个表是否存在的 SQL 语句，参数是表名(table)。
}

// RegisterDialect register a dialect to the global variable
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// Get the dialect from global variable if it exists
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
