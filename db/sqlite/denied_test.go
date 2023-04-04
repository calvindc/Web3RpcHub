package sqlite

import (
	"context"
	"reflect"
	"testing"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/refs"
)

func TestDeniedKeys_Add(t *testing.T) {
	type args struct {
		ctx     context.Context
		a       refs.FeedRef
		comment string
	}
	tests := []struct {
		name    string
		dk      DeniedKeys
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dk.Add(tt.args.ctx, tt.args.a, tt.args.comment); (err != nil) != tt.wantErr {
				t.Errorf("DeniedKeys.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeniedKeys_HasFeed(t *testing.T) {
	type args struct {
		ctx context.Context
		h   refs.FeedRef
	}
	tests := []struct {
		name string
		dk   DeniedKeys
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dk.HasFeed(tt.args.ctx, tt.args.h); got != tt.want {
				t.Errorf("DeniedKeys.HasFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeniedKeys_HasID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name string
		dk   DeniedKeys
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dk.HasID(tt.args.ctx, tt.args.id); got != tt.want {
				t.Errorf("DeniedKeys.HasID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeniedKeys_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		dk      DeniedKeys
		args    args
		want    db.ListEntry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dk.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeniedKeys.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeniedKeys.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeniedKeys_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		dk      DeniedKeys
		args    args
		want    []db.ListEntry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dk.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeniedKeys.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeniedKeys.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeniedKeys_Count(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		dk      DeniedKeys
		args    args
		want    uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dk.Count(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeniedKeys.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeniedKeys.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeniedKeys_RemoveFeed(t *testing.T) {
	type args struct {
		ctx context.Context
		r   refs.FeedRef
	}
	tests := []struct {
		name    string
		dk      DeniedKeys
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dk.RemoveFeed(tt.args.ctx, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("DeniedKeys.RemoveFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeniedKeys_RemoveID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		dk      DeniedKeys
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dk.RemoveID(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeniedKeys.RemoveID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
