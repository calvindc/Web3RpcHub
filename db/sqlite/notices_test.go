package sqlite

import (
	"context"
	"reflect"
	"testing"

	"github.com/calvindc/Web3RpcHub/db"
)

func TestPinnedNotices_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		pn      PinnedNotices
		args    args
		want    db.PinnedNotices
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pn.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PinnedNotices.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PinnedNotices.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPinnedNotices_Get(t *testing.T) {
	type args struct {
		ctx  context.Context
		name db.PinnedNoticeName
		lang string
	}
	tests := []struct {
		name    string
		pn      PinnedNotices
		args    args
		want    *db.Notice
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pn.Get(tt.args.ctx, tt.args.name, tt.args.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("PinnedNotices.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PinnedNotices.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPinnedNotices_Set(t *testing.T) {
	type args struct {
		ctx      context.Context
		name     db.PinnedNoticeName
		noticeID int64
	}
	tests := []struct {
		name    string
		pn      PinnedNotices
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pn.Set(tt.args.ctx, tt.args.name, tt.args.noticeID); (err != nil) != tt.wantErr {
				t.Errorf("PinnedNotices.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotices_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		n       Notices
		args    args
		want    db.Notice
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Notices.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Notices.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotices_RemoveID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		n       Notices
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.RemoveID(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Notices.RemoveID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotices_Save(t *testing.T) {
	type args struct {
		ctx context.Context
		p   *db.Notice
	}
	tests := []struct {
		name    string
		n       Notices
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Save(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("Notices.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
