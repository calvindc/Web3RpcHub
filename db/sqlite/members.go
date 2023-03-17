package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/db/sqlite/models"
	"github.com/calvindc/Web3RpcHub/internal/refs"
	"github.com/mattn/go-sqlite3"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// compiler assertion to ensure the struct fullfills the interface
var _ db.MembersService = (*Members)(nil)

type Members struct {
	db *sql.DB
}

// getAliases returns the JOIN-ed aliases for a member as a the wanted db type.
func (Members) getAliases(mEntry *models.Member) []db.Alias {
	if mEntry.R == nil || mEntry.R.Aliases == nil {
		return make([]db.Alias, 0)
	}
	var aliases = make([]db.Alias, len(mEntry.R.Aliases))
	for j, aEntry := range mEntry.R.Aliases {
		aliases[j].ID = aEntry.ID
		aliases[j].Feed = mEntry.PubKey.FeedRef
		aliases[j].Name = aEntry.Name
		aliases[j].Signature = aEntry.Signature
	}
	return aliases
}

func (m Members) Add(ctx context.Context, pubKey refs.FeedRef, role db.Role) (int64, error) {
	var newID int64
	err := transact(m.db, func(tx *sql.Tx) error {
		var err error
		newID, err = m.add(ctx, tx, pubKey, role)
		return err
	})
	if err != nil {
		return -1, err
	}
	return newID, nil
}

// no receiver name because it needs to use the passed transaction
func (Members) add(ctx context.Context, tx *sql.Tx, pubKey refs.FeedRef, role db.Role) (int64, error) {
	if err := role.IsValid(); err != nil {
		return -1, err
	}

	if _, err := refs.ParseFeedRef(pubKey.String()); err != nil {
		return -1, err
	}

	var newMember models.Member
	newMember.PubKey = db.DBFeedRef{FeedRef: pubKey}
	newMember.Role = int64(role)

	err := newMember.Insert(ctx, tx, boil.Infer())
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, db.ErrAlreadyAdded{Ref: pubKey}
		}
		return -1, fmt.Errorf("members: failed to insert new user: %w", err)
	}

	return newMember.ID, nil
}

func (m Members) GetByID(ctx context.Context, mid int64) (db.Member, error) {
	entry, err := models.FindMember(ctx, m.db, mid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Member{}, db.ErrNotFound
		}
		return db.Member{}, err
	}

	// resolve alias relation
	entry.R = entry.R.NewStruct()
	entry.R.Aliases, err = entry.Aliases().All(ctx, m.db)
	if err != nil {
		return db.Member{}, err
	}

	return db.Member{
		ID:      entry.ID,
		Role:    db.Role(entry.Role),
		PubKey:  entry.PubKey.FeedRef,
		Aliases: m.getAliases(entry),
	}, nil
}

// GetByFeed returns the member if it exists
func (m Members) GetByFeed(ctx context.Context, h refs.FeedRef) (db.Member, error) {
	entry, err := models.Members(qm.Where("pub_key = ?", h.String())).One(ctx, m.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Member{}, db.ErrNotFound
		}
		return db.Member{}, err
	}

	// resolve alias relation
	entry.R = entry.R.NewStruct()
	entry.R.Aliases, err = entry.Aliases().All(ctx, m.db)
	if err != nil {
		return db.Member{}, err
	}

	return db.Member{
		ID:      entry.ID,
		Role:    db.Role(entry.Role),
		PubKey:  entry.PubKey.FeedRef,
		Aliases: m.getAliases(entry),
	}, nil
}

// List returns a list of all the feeds.
func (m Members) List(ctx context.Context) ([]db.Member, error) {
	all, err := models.Members(qm.Load("Aliases")).All(ctx, m.db)
	if err != nil {
		return nil, err
	}

	var members = make([]db.Member, len(all))
	for i, entry := range all {
		members[i].ID = entry.ID
		members[i].Role = db.Role(entry.Role)
		members[i].PubKey = entry.PubKey.FeedRef
		members[i].Aliases = m.getAliases(entry)
	}

	return members, nil
}

func (m Members) Count(ctx context.Context) (uint, error) {
	count, err := models.Members().Count(ctx, m.db)
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

// RemoveFeed removes the feed from the list.
func (m Members) RemoveFeed(ctx context.Context, r refs.FeedRef) error {
	entry, err := models.Members(qm.Where("pub_key = ?", r.String())).One(ctx, m.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNotFound
		}
		return err
	}

	_, err = entry.Delete(ctx, m.db)
	if err != nil {
		return err
	}

	return nil
}

// RemoveID removes the feed from the list.
func (m Members) RemoveID(ctx context.Context, id int64) error {
	entry, err := models.FindMember(ctx, m.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNotFound
		}
		return err
	}

	_, err = entry.Delete(ctx, m.db)
	if err != nil {
		return err
	}

	return nil
}

// SetRole updates the role r of the passed memberID.
func (m Members) SetRole(ctx context.Context, id int64, r db.Role) error {
	if err := r.IsValid(); err != nil {
		return err
	}

	return transact(m.db, func(tx *sql.Tx) error {
		m, err := models.FindMember(ctx, tx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return db.ErrNotFound
			}
			return err
		}

		// find the number of other admins
		admins, err := models.Members(
			qm.Where("id != ?", id),
			qm.Where("role = ?", db.RoleAdmin),
		).Count(ctx, tx)
		if err != nil {
			return err
		}

		if admins < 1 {
			return fmt.Errorf("need at least one other admin")
		}

		m.Role = int64(r)
		_, err = m.Update(ctx, tx, boil.Whitelist("role"))
		return err
	})
}
