package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/calvindc/Web3RpcHub/internal/repository"
	"github.com/stretchr/testify/require"
)

/*func TestOpenDB(t *testing.T) {

}*/

func TestHubDB(t *testing.T) {
	testRepo := filepath.Join("testdata", t.Name())
	os.RemoveAll(testRepo)

	tr := repository.New(testRepo)

	db, err := OpenDB(tr)
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)
}

func TestOpenDB(t *testing.T) {
	type args struct {
		r repository.Interface
	}
	tests := []struct {
		name    string
		args    args
		want    *Database
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenDB(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_Close(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := Database{
				db: tt.fields.db,
			}
			if err := dt.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Database.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_transact(t *testing.T) {
	type args struct {
		db *sql.DB
		fn func(tx *sql.Tx) error
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
			if err := transact(tt.args.db, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("transact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
