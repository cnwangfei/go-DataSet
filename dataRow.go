package dataSet

import (
	"encoding/json"
	"encoding/xml"
)

type DataRow []*Data // 一行数据
// Data data
type Data struct {
	column *ColumnType // 字段类型
	value  interface{} // 值
}

// toMarshal 用于序列化输出
func (rd DataRow) toMarshal() Param {
	rowMap := make(Param, len(rd))
	for _, data := range rd {
		rowMap[data.column.Name()] = data.value
	}
	return rowMap
}

// MarshalJSON json序列化
func (rd DataRow) MarshalJSON() ([]byte, error) {
	return json.Marshal(rd.toMarshal())
}

// MarshalXML xml序列化
func (rd DataRow) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	nerRow := rd.toMarshal()
	start.Name = xml.Name{
		Space: "",
		Local: "row",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range nerRow {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
