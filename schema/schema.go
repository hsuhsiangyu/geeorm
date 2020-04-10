package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// 实现 ORM 框架中最为核心的转换——对象(object)和表(table)的转换。

// Field represents a column of database
type Field struct {
	Name string     // 字段名
	Type string		// 类型
	Tag  string		// 约束条件
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}    // 被映射的对象
	Name       string		  // 表名
	Fields     []*Field		  // 字段
	FieldNames []string
	fieldMap   map[string]*Field
}

// GetField returns field by name
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 根据数据库中列的顺序，从对象中找到对应的值，按顺序平铺。
// Values return the values of dest's member variables
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// Parse a struct to a Schema instance
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
