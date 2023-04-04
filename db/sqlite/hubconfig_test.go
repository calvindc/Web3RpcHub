package sqlite

import (
	"context"
	"reflect"
	"testing"

	"github.com/calvindc/Web3RpcHub/db"
)

func TestConfig_GetPrivacyMode(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       Config
		args    args
		want    db.PrivacyMode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetPrivacyMode(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.GetPrivacyMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.GetPrivacyMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetPrivacyMode(t *testing.T) {
	type args struct {
		ctx context.Context
		pm  db.PrivacyMode
	}
	tests := []struct {
		name    string
		c       Config
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetPrivacyMode(tt.args.ctx, tt.args.pm); (err != nil) != tt.wantErr {
				t.Errorf("Config.SetPrivacyMode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_GetDefaultLanguage(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       Config
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetDefaultLanguage(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.GetDefaultLanguage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Config.GetDefaultLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDefaultLanguage(t *testing.T) {
	type args struct {
		ctx     context.Context
		langTag string
	}
	tests := []struct {
		name    string
		c       Config
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetDefaultLanguage(tt.args.ctx, tt.args.langTag); (err != nil) != tt.wantErr {
				t.Errorf("Config.SetDefaultLanguage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
