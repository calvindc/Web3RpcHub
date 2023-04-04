package sqlite

import (
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestAuthFallback_Check(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		af      AuthFallback
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.af.Check(tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthFallback.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthFallback.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthFallback_SetPassword(t *testing.T) {
	type args struct {
		ctx      context.Context
		memberID int64
		password string
	}
	tests := []struct {
		name    string
		af      AuthFallback
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.af.SetPassword(tt.args.ctx, tt.args.memberID, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("AuthFallback.SetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthFallback_SetPasswordWithToken(t *testing.T) {
	type args struct {
		ctx        context.Context
		resetToken string
		password   string
	}
	tests := []struct {
		name    string
		af      AuthFallback
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.af.SetPasswordWithToken(tt.args.ctx, tt.args.resetToken, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("AuthFallback.SetPasswordWithToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthFallback_CreateResetToken(t *testing.T) {
	type args struct {
		ctx             context.Context
		createdByMember int64
		forMember       int64
	}
	tests := []struct {
		name    string
		af      AuthFallback
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.af.CreateResetToken(tt.args.ctx, tt.args.createdByMember, tt.args.forMember)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthFallback.CreateResetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthFallback.CreateResetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deleteConsumedResetTokens(t *testing.T) {
	type args struct {
		tx boil.ContextExecutor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteConsumedResetTokens(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("deleteConsumedResetTokens() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
