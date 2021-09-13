package dto

import "encoding/xml"
// RCurrencies struct
type RCurrencies struct {
	Item        []RCurrency `xml:"item"`
	Description string      `xml:"description"`
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	Copyright   string      `xml:"copyright"`
	Generator   string      `xml:"generator"`
	Date        string      `xml:"date"`
}
// RCurrency struct
type RCurrency struct {
	XMLName xml.Name `xml:"item"`
	TITLE   string   `xml:"fullname"`
	CODE    string   `xml:"title"`
	VALUE   float32  `xml:"description"`
}
