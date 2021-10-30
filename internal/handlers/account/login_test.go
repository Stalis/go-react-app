package account

import (
	"io"
	"strings"
	"testing"
)

func TestLoginRequest_FromJSON(t *testing.T) {
	type fields struct {
		Username     string
		PasswordHash string
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "valid JSON",
			fields:  fields{},
			args:    args{r: strings.NewReader("{\"username\": \"asd\", \"passwordHash\": \"123\" }")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoginRequest{
				Username:     tt.fields.Username,
				PasswordHash: tt.fields.PasswordHash,
			}
			if err := l.FromJSON(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("LoginRequest.FromJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
