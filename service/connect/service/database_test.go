package service

import (
	"log"
	"testing"
)

func Test_getUsername(t *testing.T) {
	InitDataBase()

	type args struct {
		uid string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid get username",
			args: args{
				uid: "1",
			},
			want:    "root",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUsername(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getUsername() got = %v, want %v", got, tt.want)
			}
			log.Println(got)
		})
	}
}
