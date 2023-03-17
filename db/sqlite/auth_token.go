package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/db/sqlite/models"
	"github.com/calvindc/Web3RpcHub/internal/randutil"
	"github.com/mattn/go-sqlite3"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// compiler assertion to ensure the struct fullfills the interface
var _ db.AuthWithTokenService = (*AuthWitToken)(nil)

type AuthWitToken struct {
	db *sql.DB
}

const protTokenLength = 32

// CreateToken is used to generate a token that is stored inside a cookie.
// It is used after a valid solution for a challenge was provided.
func (a AuthWitToken) CreateToken(ctx context.Context, memberID int64) (string, error) {

	var newToken = models.SIWSSBSession{
		MemberID: memberID,
	}

	err := transact(a.db, func(tx *sql.Tx) error {

		// check the member is registerd
		if _, err := models.FindMember(ctx, tx, newToken.MemberID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return db.ErrNotFound
			}
			return err
		}

		inserted := false
	trying: // keep trying until we inserted in unused token
		for tries := 100; tries > 0; tries-- {

			// generate an new token
			newToken.Token = randutil.String(protTokenLength)

			// insert the new token
			cols := boil.Whitelist(models.SIWSSBSessionColumns.Token, models.SIWSSBSessionColumns.MemberID)
			err := newToken.Insert(ctx, tx, cols)
			if err != nil {
				var sqlErr sqlite3.Error
				if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					// generated an existing token, retry
					continue trying
				}
				return err
			}
			inserted = true
			break // no error means it worked!
		}

		if !inserted {
			return errors.New("admindb: failed to generate a fresh token in a reasonable amount of time")
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return newToken.Token, nil
}

const sessionTimeout = time.Hour * 24

// CheckToken checks if the passed token is still valid and returns the member id if so
func (a AuthWitToken) CheckToken(ctx context.Context, token string) (int64, error) {
	var memberID int64

	err := transact(a.db, func(tx *sql.Tx) error {
		session, err := models.SIWSSBSessions(qm.Where("token = ?", token)).One(ctx, a.db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return db.ErrNotFound
			}
			return err
		}

		if time.Since(session.CreatedAt) > sessionTimeout {
			_, err = session.Delete(ctx, tx)
			if err != nil {
				return err
			}

			return errors.New("sign-in with ssb: session expired")
		}

		memberID = session.MemberID
		return nil
	})
	if err != nil {
		return -1, err
	}

	return memberID, nil
}

// RemoveToken removes a single token from the database
func (a AuthWitToken) RemoveToken(ctx context.Context, token string) error {
	_, err := models.SIWSSBSessions(qm.Where("token = ?", token)).DeleteAll(ctx, a.db)
	return err
}

// WipeTokensForMember deletes all tokens currently held for that member
func (a AuthWitToken) WipeTokensForMember(ctx context.Context, memberID int64) error {
	return transact(a.db, func(tx *sql.Tx) error {
		_, err := models.SIWSSBSessions(qm.Where("member_id = ?", memberID)).DeleteAll(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return db.ErrNotFound
			}
			return err
		}
		return nil
	})
}
