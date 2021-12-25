package dataSet

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

// DateTime 时间
// 重命名时间用于格式化输出
type DateTime struct {
	time.Time
}

// MarshalJSON json序列化
func (t DateTime) MarshalJSON() ([]byte, error) {
	new := t.Time.Format("2006-01-02 15:04:05.000")
	return json.Marshal(new)
}

// MarshalXML xml序列化
func (t DateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	new := t.Time.Format("2006-01-02 15:04:05.000")
	return e.EncodeElement(new, start)
}
