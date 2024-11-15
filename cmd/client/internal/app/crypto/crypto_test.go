package crypto

import (
	"reflect"
	"testing"
	"time"
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
			if reflect.TypeOf(got) == reflect.TypeOf((*Crypto)(nil)).Elem() {
				t.Errorf("not crypto")
			}
		})
	}
}

func TestCryptoEngine_Key(t *testing.T) {
	type fields struct {
		timestamp time.Time
	}
	type args struct {
		usr  string
		pass string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "positive test #1",
			args: args{usr: "x0o1@ya.ru", pass: "testPass"},
			want: []byte{18, 15, 39, 226, 196, 52, 136, 66, 178, 52, 90, 38, 86, 33, 125, 112, 21, 74, 160, 144, 226, 62, 89, 106, 177, 57, 186, 13, 217, 47, 117, 233},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CryptoEngine{
				timestamp: tt.fields.timestamp,
			}
			if got := c.Key(tt.args.usr, tt.args.pass); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CryptoEngine.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCryptoEngine_Enc(t *testing.T) {
	type fields struct {
		timestamp time.Time
	}
	type args struct {
		plaintext string
		key       []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "negative test #1",
			fields:  fields{timestamp: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)},
			wantErr: true,
		},
		{
			name:    "positive test #1",
			fields:  fields{timestamp: time.Now()},
			args:    args{plaintext: "t", key: []byte{18, 15, 39, 226, 196, 52, 136, 66, 178, 52, 90, 38, 86, 33, 125, 112, 21, 74, 160, 144, 226, 62, 89, 106, 177, 57, 186, 13, 217, 47, 117, 233}},
			wantErr: false,
		},
		{
			name:    "negative test #2",
			fields:  fields{timestamp: time.Now()},
			args:    args{plaintext: "t", key: []byte{18, 15, 39, 226, 196, 52, 136, 66, 178, 52, 90, 38, 86, 33, 125, 112, 21, 74, 160, 144, 226, 62, 89, 106, 177, 57, 186, 13, 217, 47, 117}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CryptoEngine{
				timestamp: tt.fields.timestamp,
			}
			_, err := c.Enc(tt.args.plaintext, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("CryptoEngine.Enc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCryptoEngine_Dec(t *testing.T) {
	type fields struct {
		timestamp time.Time
	}
	type args struct {
		ciphertext string
		key        []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "positive test #1",
			fields: fields{timestamp: time.Now()},
			args:   args{ciphertext: "aa200578091af3a9bfee69a972eafb734639d81fb78afb8a0af5866198", key: []byte{18, 15, 39, 226, 196, 52, 136, 66, 178, 52, 90, 38, 86, 33, 125, 112, 21, 74, 160, 144, 226, 62, 89, 106, 177, 57, 186, 13, 217, 47, 117, 233}},
			want:   "t",
		},
		{
			name:    "negative test #1",
			fields:  fields{timestamp: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)},
			wantErr: true,
		},
		{
			name:    "negative test #2",
			fields:  fields{timestamp: time.Now()},
			args:    args{ciphertext: "Zaa200578091af3a9bfee69a972eafb734639d81fb78afb8a0af5866198", key: []byte{18, 15, 39, 226, 196, 52, 136, 66, 178, 52, 90, 38, 86, 33, 125, 112, 21, 74, 160, 144, 226, 62, 89, 106, 177, 57, 186, 13, 217, 47, 117, 233}},
			wantErr: true,
		},
		{
			name:    "negative test #3",
			fields:  fields{timestamp: time.Now()},
			args:    args{ciphertext: "aa200578091af3a9bfee69a972eafb734639d81fb78afb8a0af5866198", key: []byte{18, 15, 39, 226, 196, 52, 136, 66, 178, 52, 90, 38, 86, 33, 125, 112, 21, 74, 160, 144, 226, 62, 89, 106, 177, 57, 186, 13, 217, 47, 117}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CryptoEngine{
				timestamp: tt.fields.timestamp,
			}
			got, err := c.Dec(tt.args.ciphertext, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("CryptoEngine.Dec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CryptoEngine.Dec() = %v, want %v", got, tt.want)
			}
		})
	}
}
