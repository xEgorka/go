// Package models contains description of service models.
package models

import "time"

// RequestAuth describes PostUserRegister and PostUserLogin handlers request.
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
