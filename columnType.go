package dataSet

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// ColumnType 列类型
type ColumnType struct {
	*sql.ColumnType
	Index uint16 `json:"index" xml:"index"`
}

// columnTypeMarshal 序列化用的列类型
type columnTypeMarshal struct {
	Index     uint16 `json:"index" xml:"index"`
	Name      string `json:"name" xml:"name"`
	typeName  string
	length    int64  // 字符串长度
	precision int64  // 数字精度
	scale     int64  // 小数位数
	TypeName  string `json:"type_name" xml:"type_name"`
	IsNull    bool   `json:"is_null" xml:"is_null"`
}

func (col *ColumnType) toMarshal() (c *columnTypeMarshal) {
	c = new(columnTypeMarshal)
	c.Index = col.Index
	c.Name = col.Name()
	c.typeName = col.DatabaseTypeName()
	// TypeName
	var typeExt string
	l, ok := col.Length()
	if ok {
		c.length = l
		typeExt = fmt.Sprintf("(%v)", l)
	}
	l1, l2, ok := col.DecimalSize()
	if ok {
		c.precision = l1
		c.scale = l2
		typeExt = fmt.Sprintf("(%v,%v)", l1, l2)
	}
	c.TypeName = fmt.Sprintf("%v%v", col.DatabaseTypeName(), typeExt)
	// isnull
	isnull, ok := col.Nullable()
	if !ok {
		isnull = false
	}
	c.IsNull = isnull
	return
}

// MarshalXML xml序列化
func (col *ColumnType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "col",
	}
	return e.EncodeElement(col.toMarshal(), start)
}

// MarshalJSON json序列化
func (col *ColumnType) MarshalJSON() ([]byte, error) {
	return json.Marshal(col.toMarshal())
}

func (col *ColumnType) String() string {
	c := col.toMarshal()
	return fmt.Sprintf("Index: %v Name: %v TypeName: %v IsNull: %v", c.Index, c.Name, c.TypeName, c.IsNull)
}
