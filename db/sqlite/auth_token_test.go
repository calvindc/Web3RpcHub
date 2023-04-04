package sqlite

import (
	"context"
	"testing"
)

func TestAuthWitToken_CreateToken(t *testing.T) {
	type args struct {
		ctx      context.Context
		memberID int64
	}
	tests := []struct {
		name    string
		a       AuthWitToken
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.CreateToken(tt.args.ctx, tt.args.memberID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthWitToken.CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthWitToken.CreateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthWitToken_CheckToken(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		a       AuthWitToken
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.CheckToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthWitToken.CheckToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthWitToken.CheckToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthWitToken_RemoveToken(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		a       AuthWitToken
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.RemoveToken(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("AuthWitToken.RemoveToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthWitToken_WipeTokensForMember(t *testing.T) {
	type args struct {
		ctx      context.Context
		memberID int64
	}
	tests := []struct {
		name    string
		a       AuthWitToken
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.WipeTokensForMember(tt.args.ctx, tt.args.memberID); (err != nil) != tt.wantErr {
				t.Errorf("AuthWitToken.WipeTokensForMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
