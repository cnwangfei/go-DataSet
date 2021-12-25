package dataSet

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"testing"
)

func TestNewDataSetFromRows(t *testing.T) {
	connStr := fmt.Sprintf("server=%s\\%s;user id=%s;password=%s;database=%s;app name=PivasApi;",
		"127.0.0.1", "SQL2014", "sa", "`1q", "master")
	db, err := gorm.Open(sqlserver.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sql := `
CREATE  TABLE #test
(
	a NUMERIC(18,4)
	,b DECIMAL(18,4)
	,c INT
	,e BIT
	,f VARCHAR(100)
	,g DATETIME
	,H VARCHAR(100)
	,i FLOAT
)

INSERT INTO dbo.#test 
VALUES	( 1.0021,1.003,2,0,'aaa','2021-12-20 12:13:14.111','<>','3.141592654')
,( 1.0020,1.00356,123456,1,'bbb','2021-12-21 01:10:14.101','<>','3.141592654')
,( 1.0020,1.00356,123456,1,'中文','2021-12-22 04:59:59.999','!@#','3.141592654')

SELECT * FROM #test

DROP TABLE #test
`
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return
	}
	ds, err := NewDataSetFromRows(rows)
	if err != nil {
		panic(err)
	}

	h := ds
	// dataset内容
	fmt.Println(h)
	fmt.Println("========================================================")
	// xml序列化
	x, err := xml.Marshal(h)
	if err != nil {
		panic(err)
	}
	fmt.Println("xml序列化")
	fmt.Printf("%s\n", x)
	fmt.Println("========================================================")
	// json序列化
	j, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	fmt.Println("json序列化")
	fmt.Printf("%s\n", j)
}
