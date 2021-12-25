package dataSet

import "encoding/xml"

type Param map[string]interface{}

// MarshalXML xml序列化
func (p Param) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "params",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range p {
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
