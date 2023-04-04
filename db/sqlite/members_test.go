package sqlite

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/db/sqlite/models"
	"github.com/calvindc/Web3RpcHub/refs"
)

func TestMembers_getAliases(t *testing.T) {
	type args struct {
		mEntry *models.Member
	}
	tests := []struct {
		name string
		m    Members
		args args
		want []db.Alias
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getAliases(tt.args.mEntry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Members.getAliases() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembers_Add(t *testing.T) {
	type args struct {
		ctx    context.Context
		pubKey refs.FeedRef
		role   db.Role
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Add(tt.args.ctx, tt.args.pubKey, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("Members.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Members.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembers_add(t *testing.T) {
	type args struct {
		ctx    context.Context
		tx     *sql.Tx
		pubKey refs.FeedRef
		role   db.Role
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.add(tt.args.ctx, tt.args.tx, tt.args.pubKey, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("Members.add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Members.add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembers_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		mid int64
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		want    db.Member
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetByID(tt.args.ctx, tt.args.mid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Members.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Members.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembers_GetByFeed(t *testing.T) {
	type args struct {
		ctx context.Context
		h   refs.FeedRef
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		want    db.Member
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetByFeed(tt.args.ctx, tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("Members.GetByFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Members.GetByFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembers_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		want    []db.Member
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Members.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Members.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembers_Count(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		want    uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Count(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Members.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Members.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembers_RemoveFeed(t *testing.T) {
	type args struct {
		ctx context.Context
		r   refs.FeedRef
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.RemoveFeed(tt.args.ctx, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Members.RemoveFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMembers_RemoveID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.RemoveID(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Members.RemoveID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMembers_SetRole(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
		r   db.Role
	}
	tests := []struct {
		name    string
		m       Members
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetRole(tt.args.ctx, tt.args.id, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Members.SetRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
