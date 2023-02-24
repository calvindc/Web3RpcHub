package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/calvindc/Web3RpcHub/internal/repository"
	migrate "github.com/rubenv/sql-migrate"
)

/*var migrationSource = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: migrations,
	Root:       "migrations",
}*/

var migrationSource = &migrate.FileMigrationSource{
	Dir: "db/sqlite/migrations",
}

type Database struct {
	db *sql.DB
}

func OpenDB(r repository.GetPathInterface) (*Database, error) {
	dbFileName := r.GetPath("hubdb")

	if dir := filepath.Dir(dbFileName); dir != "" {
		err := os.MkdirAll(dir, 0700)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("db: failed to create folder for database (%q): %w", dir, err)
		}
	}
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, fmt.Errorf("db: failed to open sqlite database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: sqlite ping failed: %w", err)
	}
	n, err := migrate.Exec(db, "sqlite3", migrationSource, migrate.Up)
	if err != nil {
		return nil, fmt.Errorf("db: failed to apply database migration: %w", err)
	}
	fmt.Printf("db: Applied %d migrations!\n", n)

	//清除旧的邀请码，重置tokens
	/*go func() {
		senvenDays := time.Hour * 24 * 7
		ticker := time.NewTicker(senvenDays)
		for range ticker.C {
			err := transact(db, func(tx *sql.Tx) error {
				if err := deleteConsumedResetTokens(tx); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				fmt.Printf("db: failed to clean up old invites: %s", err)
			}
		}
	}()*/

	return &Database{}, nil
}

func transact(db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return fmt.Errorf("transact: could not begin transaction: %w", err)
	}

	if err = fn(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			err = fmt.Errorf("rollback failed after %s: %s", err.Error(), err2.Error())
		} else {
			err = fmt.Errorf("transaction failed, and rolling back : %w", err)
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("transact: could not commit transaction: %w", err)
	}

	return nil
}

// since reset tokens are marked as invalid so that the code can't be generated twice,
// they need to be deleted periodically.
/*func deleteConsumedResetTokens(tx boil.ContextExecutor) error {
	_, err := models.FallbackResetTokens(qm.Where("active = false")).DeleteAll(context.Background(), tx)
	if err != nil {
		return fmt.Errorf("admindb: failed to delete used reset tokens: %w", err)
	}
	return nil
}
*/
