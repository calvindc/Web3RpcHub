package sqlite

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/refs"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func TestAliases_Register(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx       context.Context
		alias     string
		userFeed  refs.FeedRef
		signature []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Aliases{
				db: tt.fields.db,
			}
			if err := a.Register(tt.args.ctx, tt.args.alias, tt.args.userFeed, tt.args.signature); (err != nil) != tt.wantErr {
				t.Errorf("Aliases.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliases_List(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []db.Alias
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Aliases{
				db: tt.fields.db,
			}
			got, err := a.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aliases.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aliases.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliases_Revoke(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx   context.Context
		alias string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Aliases{
				db: tt.fields.db,
			}
			if err := a.Revoke(tt.args.ctx, tt.args.alias); (err != nil) != tt.wantErr {
				t.Errorf("Aliases.Revoke() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliases_Resolve(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    db.Alias
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Aliases{
				db: tt.fields.db,
			}
			got, err := a.Resolve(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aliases.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aliases.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliases_GetByID(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    db.Alias
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Aliases{
				db: tt.fields.db,
			}
			got, err := a.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aliases.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aliases.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliases_findOne(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		by  qm.QueryMod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    db.Alias
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Aliases{
				db: tt.fields.db,
			}
			got, err := a.findOne(tt.args.ctx, tt.args.by)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aliases.findOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aliases.findOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
