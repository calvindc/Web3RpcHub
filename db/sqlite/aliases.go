package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/db/sqlite/models"
	"github.com/calvindc/Web3RpcHub/internal/refs"
	"github.com/mattn/go-sqlite3"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// compiler 断言，确保结构完全填充AliasesDBService接口内的方法
var _ db.AliasesDBService = (*Aliases)(nil)

type Aliases struct {
	db *sql.DB
}

// Register 客户端通过验证后，注册alias(alias,userFeed,signature)
func (a Aliases) Register(ctx context.Context, alias string, userFeed refs.FeedRef, signature []byte) error {
	return transact(a.db, func(tx *sql.Tx) error {
		// 在member中查找userFeed
		memberEntry, err := models.Members(qm.Where("pub_key = ?", userFeed.String())).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return db.ErrNotFound
			}
			return err
		}

		var newEntry models.Alias
		newEntry.Name = alias
		newEntry.MemberID = memberEntry.ID
		newEntry.Signature = signature

		err = newEntry.Insert(ctx, tx, boil.Infer())
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return db.ErrAliasTaken{Name: alias}
		}
		return err
	})
}

// List 查询所有注册的aliases
func (a Aliases) List(ctx context.Context) ([]db.Alias, error) {
	all, err := models.Aliases(qm.Load("Member")).All(ctx, a.db)
	if err != nil {
		return nil, err
	}

	var aliases = make([]db.Alias, len(all))
	for i, entry := range all {
		aliases[i] = db.Alias{
			ID:        entry.ID,
			Name:      entry.Name,
			Feed:      entry.R.Member.PubKey.FeedRef,
			Signature: entry.Signature,
		}
	}

	return aliases, nil
}

// Revoke 删除某个alias
func (a Aliases) Revoke(ctx context.Context, alias string) error {
	return transact(a.db, func(tx *sql.Tx) error {
		qry := append([]qm.QueryMod{qm.Load("Member")}, qm.Where("name = ?", alias))

		entry, err := models.Aliases(qry...).One(ctx, a.db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return db.ErrNotFound
			}
			return err
		}

		_, err = entry.Delete(ctx, tx)
		return err
	})
}

// Resolve 返回该alias的相关信息，不存在则返回错误
func (a Aliases) Resolve(ctx context.Context, name string) (db.Alias, error) {
	return a.findOne(ctx, qm.Where("name = ?", name))
}

// GetByID returns the alias for that ID or an error
func (a Aliases) GetByID(ctx context.Context, id int64) (db.Alias, error) {
	return a.findOne(ctx, qm.Where("id = ?", id))
}

func (a Aliases) findOne(ctx context.Context, by qm.QueryMod) (db.Alias, error) {
	var found db.Alias

	// 构造查询，该查询解析Member关系
	qry := append([]qm.QueryMod{qm.Load("Member")}, by)
	entry, err := models.Aliases(qry...).One(ctx, a.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return found, db.ErrNotFound
		}
		return found, err
	}

	// unpack models into roomdb type
	found.ID = entry.ID
	found.Name = entry.Name
	found.Signature = entry.Signature
	found.Feed = entry.R.Member.PubKey.FeedRef

	return found, nil
}
