package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/db/sqlite/models"
	"github.com/calvindc/Web3RpcHub/refs"
	"github.com/mattn/go-sqlite3"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// compiler assertion to ensure the struct fullfills the interface
var _ db.DeniedKeysService = (*DeniedKeys)(nil)

// The DeniedKeys is backed by the members table
type DeniedKeys struct {
	db *sql.DB
}

// Add adds the feed to the list.
func (dk DeniedKeys) Add(ctx context.Context, a refs.FeedRef, comment string) error {
	// TODO: better valid
	if _, err := refs.ParseFeedRef(a.String()); err != nil {
		return err
	}

	var entry models.DeniedKey
	entry.PubKey.FeedRef = a
	entry.Comment = comment

	err := entry.Insert(ctx, dk.db, boil.Whitelist("pub_key", "comment"))
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return db.ErrAlreadyAdded{Ref: a}
		}

		return fmt.Errorf("Denied-list: failed to insert new entry %s: %w - type:%T", entry.PubKey, err, err)
	}

	return nil
}

// HasFeed returns true if a feed is on the list.
func (dk DeniedKeys) HasFeed(ctx context.Context, h refs.FeedRef) bool {
	_, err := models.DeniedKeys(qm.Where("pub_key = ?", h.String())).One(ctx, dk.db)
	if err != nil {
		return false
	}
	return true
}

// HasID returns true if a feed is on the list.
func (dk DeniedKeys) HasID(ctx context.Context, id int64) bool {
	_, err := models.FindDeniedKey(ctx, dk.db, id)
	if err != nil {
		return false
	}
	return true
}

// GetByID returns the entry if a feed with that ID is on the list.
func (dk DeniedKeys) GetByID(ctx context.Context, id int64) (db.ListEntry, error) {
	var entry db.ListEntry
	found, err := models.FindDeniedKey(ctx, dk.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entry, db.ErrNotFound
		}
		return entry, err
	}

	entry.ID = found.ID
	entry.PubKey = found.PubKey.FeedRef
	entry.Comment = found.Comment
	entry.CreatedAt = found.CreatedAt
	return entry, nil
}

// List returns a list of all the feeds.
func (dk DeniedKeys) List(ctx context.Context) ([]db.ListEntry, error) {
	all, err := models.DeniedKeys().All(ctx, dk.db)
	if err != nil {
		return nil, err
	}
	n := len(all)

	var lst = make([]db.ListEntry, n)
	for i, entry := range all {
		lst[i].ID = entry.ID
		lst[i].PubKey = entry.PubKey.FeedRef
		lst[i].Comment = entry.Comment
		lst[i].CreatedAt = entry.CreatedAt
	}

	return lst, nil
}

func (dk DeniedKeys) Count(ctx context.Context) (uint, error) {
	count, err := models.DeniedKeys().Count(ctx, dk.db)
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

// RemoveFeed removes the feed from the list.
func (dk DeniedKeys) RemoveFeed(ctx context.Context, r refs.FeedRef) error {
	entry, err := models.DeniedKeys(qm.Where("pub_key = ?", r.String())).One(ctx, dk.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNotFound
		}
		return err
	}

	_, err = entry.Delete(ctx, dk.db)
	if err != nil {
		return err
	}

	return nil
}

// RemoveID removes the feed from the list.
func (dk DeniedKeys) RemoveID(ctx context.Context, id int64) error {
	entry, err := models.FindDeniedKey(ctx, dk.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNotFound
		}
		return err
	}

	_, err = entry.Delete(ctx, dk.db)
	if err != nil {
		return err
	}

	return nil
}
