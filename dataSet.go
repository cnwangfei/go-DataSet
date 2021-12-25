package dataSet

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"reflect"
)

type DataSet []*DataTable

// NewDataSetFromRows 从sql.Rows生成DataSet
func NewDataSetFromRows(sqlRows *sql.Rows) (ds DataSet, err error) {
	flag := true // 可以处理多个数据库查询结果集
	for flag {
		dt, err := newDataTableFromRows(sqlRows)
		if err != nil {
			return nil, err
		}
		ds = append(ds, dt)
		flag = sqlRows.NextResultSet() //检查是否还有其他结果集
	}
	return
}

// MarshalXML xml序列化
func (ds DataSet) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	types := reflect.TypeOf(ds)
	start.Name = xml.Name{
		Space: "",
		Local: types.Name(),
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for i, dt := range ds {
		dtStart := xml.StartElement{
			Name: xml.Name{
				Space: "",
				Local: fmt.Sprintf("%v%v", reflect.TypeOf(*dt).Name(), i),
			},
			Attr: []xml.Attr{},
		}
		if err := dt.MarshalXML(e, dtStart); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (ds DataSet) String() string {
	var reStr string
	for i, dt := range ds {
		reStr = fmt.Sprintf("%v\nTable%v\n%v", reStr, i, dt.String())
	}
	return reStr
}
