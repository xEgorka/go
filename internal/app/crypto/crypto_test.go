package crypto

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Crypto
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if reflect.TypeOf(got) == reflect.TypeOf((*CryptoEngine)(nil)).Elem() {
				t.Errorf("not crypto engine")
			}
		})
	}
}

func TestCryptoEngine_Hash(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name    string
		c       *CryptoEngine
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{pass: "testPass"},
			wantErr: false,
		},
		{
			name:    "negative test #1",
			args:    args{pass: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CryptoEngine{}
			_, err := c.Hash(tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("CryptoEngine.Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCryptoEngine_Verify(t *testing.T) {
	type args struct {
		pass string
		hash string
	}
	tests := []struct {
		name    string
		c       *CryptoEngine
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{pass: "testPass", hash: "$2a$10$iYPAVEHMkzEQgG/fNB.M8eFO8sZ/6DLWiNvcsKg.atUSC0oXSZAsO"},
			wantErr: false,
		},
		{
			name:    "negative test #1",
			args:    args{pass: "testPass", hash: "badhash"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CryptoEngine{}
			if err := c.Verify(tt.args.pass, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("CryptoEngine.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
