package config

import "testing"

func TestSetup(t *testing.T) {
	tests := []struct {
		want *Config
		name string
	}{
		{
			name: "positive test #1",
			want: &Config{
				AppAddr: "https://localhost:8080/",
				DBPath:  "gophkeeper.db"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := Setup()
			if c.AppAddr != tt.want.AppAddr {
				t.Errorf("Setup() returns bad config")
			}
			if c.DBPath != tt.want.DBPath {
				t.Errorf("Setup() returns bad config")
			}
		})
	}
}
