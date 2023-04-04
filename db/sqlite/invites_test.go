package sqlite

import (
	"context"
	"reflect"
	"testing"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/refs"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestInvites_Create(t *testing.T) {
	type args struct {
		ctx       context.Context
		createdBy int64
	}
	tests := []struct {
		name    string
		i       Invites
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Create(tt.args.ctx, tt.args.createdBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invites.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Invites.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvites_Consume(t *testing.T) {
	type args struct {
		ctx       context.Context
		token     string
		newMember refs.FeedRef
	}
	tests := []struct {
		name    string
		i       Invites
		args    args
		want    db.Invite
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Consume(tt.args.ctx, tt.args.token, tt.args.newMember)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invites.Consume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invites.Consume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deleteConsumedInvites(t *testing.T) {
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
			if err := deleteConsumedInvites(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("deleteConsumedInvites() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInvites_GetByToken(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		i       Invites
		args    args
		want    db.Invite
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.GetByToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invites.GetByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invites.GetByToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvites_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		i       Invites
		args    args
		want    db.Invite
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invites.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invites.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvites_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		i       Invites
		args    args
		want    []db.Invite
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invites.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invites.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvites_Count(t *testing.T) {
	type args struct {
		ctx        context.Context
		onlyActive bool
	}
	tests := []struct {
		name    string
		i       Invites
		args    args
		want    uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Count(tt.args.ctx, tt.args.onlyActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invites.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Invites.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvites_Revoke(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		i       Invites
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.Revoke(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Invites.Revoke() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getHashedToken(t *testing.T) {
	type args struct {
		b64tok string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getHashedToken(tt.args.b64tok)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHashedToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getHashedToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
