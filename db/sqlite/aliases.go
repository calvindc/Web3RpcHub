package sqlite

import (
	"database/sql"
)

// compiler 断言，确保结构完全填充接口
//var _ db.AliasesDBService = (*Aliases)(nil)

type Aliases struct {
	db *sql.DB
}

// Resolve returns all the relevant information for that alias or an error if it doesnt exist
/*func (a Aliases) Resolve(ctx context.Context, name string) (db.Alias, error) {
	return a.findOne(ctx, qm.Where("name = ?", name))
}
*/
/*func (a Aliases) findOne(ctx context.Context, by qm.QueryMod) (db.Alias, error) {
	var found db.Alias

	// construct query which resolves the Member relation and by which we shoudl look for it
	qry := append([]qm.QueryMod{qm.Load("Member")}, by)

	entry, err := Aliases(qry...).One(ctx, a.db)
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
}*/
