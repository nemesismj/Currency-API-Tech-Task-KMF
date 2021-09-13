package model

import "time"
// RCurrency struct
type RCurrency struct {
	ID     int       `json:"id,omitempty"`
	TITLE  string    `json:"title,omitempty"`
	CODE   string    `json:"code,omitempty"`
	VALUE  float32   `json:"value,omitempty"`
	A_DATE time.Time `json:"a_date,omitempty"`
}
