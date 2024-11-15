package auth

import (
	"testing"
)

func TestBuildJWTString(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				userID: "x0o1@ya.ru",
			},
			want: 135,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildJWTString(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildJWTString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("BuildJWTString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "negative test #1",
			args: args{
				tokenString: "hhhtttt02222",
			},
		},
		{
			name: "positive test #1",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIwMzM3MjQ1MzUsIlVzZXJJRCI6ImVlQHJyLnJ1In0.4aFlz_lViChwCqj6_eVv1j9AGU0Y1NVgdyOoO0ChGhs",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetUserID(tt.args.tokenString)
			if err == nil && tt.name == "negative test #1" {
				t.Errorf("GetUserID() no error")
				return
			}
		})
	}
}
