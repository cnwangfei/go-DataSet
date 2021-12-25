package dataSet

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type DataTable struct {
	Columns []*ColumnType `json:"columns" xml:"columns"` // 列
	Rows    []DataRow     `json:"rows" xml:"rows"`       // 行数据
}

// 从sql.Rows生成DataTable
func newDataTableFromRows(sqlRows *sql.Rows) (dt *DataTable, err error) {
	cols, err := newColumnsFromRows(sqlRows)
	if err != nil {
		return
	}
	rows, err := newRowsDataFromRows(sqlRows, cols)
	if err != nil {
		return
	}
	dt = new(DataTable)
	dt.Columns = cols
	dt.Rows = rows
	return
}

// 从sql.Rows取出字段信息
func newColumnsFromRows(sqlRows *sql.Rows) (cols []*ColumnType, err error) {
	cts, err := sqlRows.ColumnTypes()
	if err != nil {
		return
	}
	cols = make([]*ColumnType, 0, len(cts))
	for index, value := range cts {
		tempCol := new(ColumnType)
		tempCol.ColumnType = value
		tempCol.Index = uint16(index)
		cols = append(cols, tempCol)
	}
	return
}

// 从sql.Rows取出所有行数据
func newRowsDataFromRows(sqlRows *sql.Rows, cols []*ColumnType) (rows []DataRow, err error) {
	colCount := len(cols)                     // 多少字段
	scanArgs := make([]interface{}, colCount) // 用于扫描每行数据时作为字段参数
	values := make([]interface{}, colCount)   // 用于扫描每行数据时临时存储 每个字段的值
	for i := range values {
		scanArgs[i] = &values[i] // 把字段参数和值指针关联起来
	}
	for sqlRows.Next() {
		err = sqlRows.Scan(scanArgs...)
		if err != nil {
			return
		}
		row := make(DataRow, colCount)     // 新建空行数据
		for index, value := range values { // 循环把取到的值都存到新的空行数据里
			//row[index] = value
			switch value.(type) {
			case []byte:
				s := string(value.([]byte))
				i2, err := strconv.Atoi(s)
				if err == nil {
					row[index] = &Data{
						column: cols[index],
						value:  i2,
					}
					continue
				}
				f, err := strconv.ParseFloat(s, 64)
				if err == nil {
					row[index] = &Data{
						column: cols[index],
						value:  f,
					}
					continue
				}
				row[index] = &Data{
					column: cols[index],
					value:  value,
				}
			case time.Time:
				row[index] = &Data{
					column: cols[index],
					value:  DateTime{value.(time.Time)},
				}
			default:
				row[index] = &Data{
					column: cols[index],
					value:  value,
				}
			}

		}
		rows = append(rows, row)
	}
	return
}

// MarshalXML 将DataTable转为xml
func (dt DataTable) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	// 序列化列
	types := reflect.TypeOf(dt)
	colsStart := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: types.Field(0).Tag.Get("xml"),
		},
		Attr: []xml.Attr{},
	}
	if err := e.EncodeToken(colsStart); err != nil {
		return err
	}

	if err := e.Encode(dt.Columns); err != nil {
		return err
	}
	e.EncodeToken(xml.EndElement{Name: colsStart.Name})

	// 序列化行数据
	rowsStart := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: types.Field(1).Tag.Get("xml"),
		},
		Attr: []xml.Attr{},
	}
	if err := e.EncodeToken(rowsStart); err != nil {
		return err
	}

	if err := e.Encode(dt.Rows); err != nil {
		return err
	}
	e.EncodeToken(xml.EndElement{Name: rowsStart.Name})

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (dt DataTable) String() string {
	var reStr string
	reStr = "Columns:"
	for _, col := range dt.Columns {
		reStr = fmt.Sprintf("%v\n\t%v", reStr, col.String())
	}
	reStr = fmt.Sprintf("%v\nRows:", reStr)
	for rowId, row := range dt.Rows {
		reStr = fmt.Sprintf("%v\n[%v]", reStr, rowId)
		for colIndex, data := range row {
			col := dt.Columns[colIndex].toMarshal() // 对应的列
			reStr = fmt.Sprintf("%v\t%v: ", reStr, col.Name)
			reStr = fmt.Sprintf("%v%v", reStr, data.value)
		}
	}
	return reStr
}
