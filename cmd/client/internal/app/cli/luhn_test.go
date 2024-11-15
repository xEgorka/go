package cli

import (
	"testing"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !valid(4111111111111111) {
				t.Errorf("valid card invalid")
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if checksum(5555555555554444/10) != 6 {
				t.Errorf("checksum incorrect")
			}
		})
	}
}
