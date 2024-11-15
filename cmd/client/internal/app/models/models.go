// Package models contains description of service models.
package models

import "time"

// RequestAuth describes user register and login requests.
type RequestAuth struct {
	Usr  string `json:"usr"`
	Pass string `json:"pass"`
}

// UserData describes user data.
type UserData struct {
	Updated time.Time `json:"updated,omitempty"`
	Merged  time.Time `json:"merged,omitempty"`
	UserID  string    `json:"usr,omitempty"`
	ID      string    `json:"id"`
	Data    string    `json:"data,omitempty"`
	Type    int       `json:"type,omitempty"`
}

// Pass describes credentials data type.
type Pass struct {
	Website string `json:"website"`
	Login   string `json:"login"`
	Pass    string `json:"pass"`
	Note    string `json:"note,omitempty"`
}

// Text describes text data type.
type Text struct {
	Text string `json:"data"`
	Note string `json:"note,omitempty"`
}

// File describes binary data type.
type File struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Note string `json:"note,omitempty"`
}

// Card describes bank card data type.
type Card struct {
	Bank string `json:"bank,omitempty"`
	Note string `json:"note,omitempty"`
	Num  int    `json:"num"`
	CVV  int    `json:"cvv"`
	Exp  int    `json:"exp"`
}
