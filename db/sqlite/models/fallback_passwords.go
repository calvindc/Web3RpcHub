// Code generated by SQLBoiler 4.14.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// FallbackPassword is an object representing the database table.
type FallbackPassword struct {
	ID           int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	PasswordHash []byte `boil:"password_hash" json:"password_hash" toml:"password_hash" yaml:"password_hash"`
	MemberID     int64  `boil:"member_id" json:"member_id" toml:"member_id" yaml:"member_id"`

	R *fallbackPasswordR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L fallbackPasswordL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var FallbackPasswordColumns = struct {
	ID           string
	PasswordHash string
	MemberID     string
}{
	ID:           "id",
	PasswordHash: "password_hash",
	MemberID:     "member_id",
}

var FallbackPasswordTableColumns = struct {
	ID           string
	PasswordHash string
	MemberID     string
}{
	ID:           "fallback_passwords.id",
	PasswordHash: "fallback_passwords.password_hash",
	MemberID:     "fallback_passwords.member_id",
}

// Generated where

var FallbackPasswordWhere = struct {
	ID           whereHelperint64
	PasswordHash whereHelper__byte
	MemberID     whereHelperint64
}{
	ID:           whereHelperint64{field: "\"fallback_passwords\".\"id\""},
	PasswordHash: whereHelper__byte{field: "\"fallback_passwords\".\"password_hash\""},
	MemberID:     whereHelperint64{field: "\"fallback_passwords\".\"member_id\""},
}

// FallbackPasswordRels is where relationship names are stored.
var FallbackPasswordRels = struct {
	Member string
}{
	Member: "Member",
}

// fallbackPasswordR is where relationships are stored.
type fallbackPasswordR struct {
	Member *Member `boil:"Member" json:"Member" toml:"Member" yaml:"Member"`
}

// NewStruct creates a new relationship struct
func (*fallbackPasswordR) NewStruct() *fallbackPasswordR {
	return &fallbackPasswordR{}
}

func (r *fallbackPasswordR) GetMember() *Member {
	if r == nil {
		return nil
	}
	return r.Member
}

// fallbackPasswordL is where Load methods for each relationship are stored.
type fallbackPasswordL struct{}

var (
	fallbackPasswordAllColumns            = []string{"id", "password_hash", "member_id"}
	fallbackPasswordColumnsWithoutDefault = []string{}
	fallbackPasswordColumnsWithDefault    = []string{"id", "password_hash", "member_id"}
	fallbackPasswordPrimaryKeyColumns     = []string{"id"}
	fallbackPasswordGeneratedColumns      = []string{}
)

type (
	// FallbackPasswordSlice is an alias for a slice of pointers to FallbackPassword.
	// This should almost always be used instead of []FallbackPassword.
	FallbackPasswordSlice []*FallbackPassword
	// FallbackPasswordHook is the signature for custom FallbackPassword hook methods
	FallbackPasswordHook func(context.Context, boil.ContextExecutor, *FallbackPassword) error

	fallbackPasswordQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	fallbackPasswordType                 = reflect.TypeOf(&FallbackPassword{})
	fallbackPasswordMapping              = queries.MakeStructMapping(fallbackPasswordType)
	fallbackPasswordPrimaryKeyMapping, _ = queries.BindMapping(fallbackPasswordType, fallbackPasswordMapping, fallbackPasswordPrimaryKeyColumns)
	fallbackPasswordInsertCacheMut       sync.RWMutex
	fallbackPasswordInsertCache          = make(map[string]insertCache)
	fallbackPasswordUpdateCacheMut       sync.RWMutex
	fallbackPasswordUpdateCache          = make(map[string]updateCache)
	fallbackPasswordUpsertCacheMut       sync.RWMutex
	fallbackPasswordUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var fallbackPasswordAfterSelectHooks []FallbackPasswordHook

var fallbackPasswordBeforeInsertHooks []FallbackPasswordHook
var fallbackPasswordAfterInsertHooks []FallbackPasswordHook

var fallbackPasswordBeforeUpdateHooks []FallbackPasswordHook
var fallbackPasswordAfterUpdateHooks []FallbackPasswordHook

var fallbackPasswordBeforeDeleteHooks []FallbackPasswordHook
var fallbackPasswordAfterDeleteHooks []FallbackPasswordHook

var fallbackPasswordBeforeUpsertHooks []FallbackPasswordHook
var fallbackPasswordAfterUpsertHooks []FallbackPasswordHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *FallbackPassword) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *FallbackPassword) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *FallbackPassword) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *FallbackPassword) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *FallbackPassword) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *FallbackPassword) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *FallbackPassword) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *FallbackPassword) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *FallbackPassword) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackPasswordAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddFallbackPasswordHook registers your hook function for all future operations.
func AddFallbackPasswordHook(hookPoint boil.HookPoint, fallbackPasswordHook FallbackPasswordHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		fallbackPasswordAfterSelectHooks = append(fallbackPasswordAfterSelectHooks, fallbackPasswordHook)
	case boil.BeforeInsertHook:
		fallbackPasswordBeforeInsertHooks = append(fallbackPasswordBeforeInsertHooks, fallbackPasswordHook)
	case boil.AfterInsertHook:
		fallbackPasswordAfterInsertHooks = append(fallbackPasswordAfterInsertHooks, fallbackPasswordHook)
	case boil.BeforeUpdateHook:
		fallbackPasswordBeforeUpdateHooks = append(fallbackPasswordBeforeUpdateHooks, fallbackPasswordHook)
	case boil.AfterUpdateHook:
		fallbackPasswordAfterUpdateHooks = append(fallbackPasswordAfterUpdateHooks, fallbackPasswordHook)
	case boil.BeforeDeleteHook:
		fallbackPasswordBeforeDeleteHooks = append(fallbackPasswordBeforeDeleteHooks, fallbackPasswordHook)
	case boil.AfterDeleteHook:
		fallbackPasswordAfterDeleteHooks = append(fallbackPasswordAfterDeleteHooks, fallbackPasswordHook)
	case boil.BeforeUpsertHook:
		fallbackPasswordBeforeUpsertHooks = append(fallbackPasswordBeforeUpsertHooks, fallbackPasswordHook)
	case boil.AfterUpsertHook:
		fallbackPasswordAfterUpsertHooks = append(fallbackPasswordAfterUpsertHooks, fallbackPasswordHook)
	}
}

// One returns a single fallbackPassword record from the query.
func (q fallbackPasswordQuery) One(ctx context.Context, exec boil.ContextExecutor) (*FallbackPassword, error) {
	o := &FallbackPassword{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for fallback_passwords")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all FallbackPassword records from the query.
func (q fallbackPasswordQuery) All(ctx context.Context, exec boil.ContextExecutor) (FallbackPasswordSlice, error) {
	var o []*FallbackPassword

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to FallbackPassword slice")
	}

	if len(fallbackPasswordAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all FallbackPassword records in the query.
func (q fallbackPasswordQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count fallback_passwords rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q fallbackPasswordQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if fallback_passwords exists")
	}

	return count > 0, nil
}

// Member pointed to by the foreign key.
func (o *FallbackPassword) Member(mods ...qm.QueryMod) memberQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.MemberID),
	}

	queryMods = append(queryMods, mods...)

	return Members(queryMods...)
}

// LoadMember allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (fallbackPasswordL) LoadMember(ctx context.Context, e boil.ContextExecutor, singular bool, maybeFallbackPassword interface{}, mods queries.Applicator) error {
	var slice []*FallbackPassword
	var object *FallbackPassword

	if singular {
		var ok bool
		object, ok = maybeFallbackPassword.(*FallbackPassword)
		if !ok {
			object = new(FallbackPassword)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeFallbackPassword)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeFallbackPassword))
			}
		}
	} else {
		s, ok := maybeFallbackPassword.(*[]*FallbackPassword)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeFallbackPassword)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeFallbackPassword))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &fallbackPasswordR{}
		}
		args = append(args, object.MemberID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &fallbackPasswordR{}
			}

			for _, a := range args {
				if a == obj.MemberID {
					continue Outer
				}
			}

			args = append(args, obj.MemberID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`members`),
		qm.WhereIn(`members.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Member")
	}

	var resultSlice []*Member
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Member")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for members")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for members")
	}

	if len(memberAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Member = foreign
		if foreign.R == nil {
			foreign.R = &memberR{}
		}
		foreign.R.FallbackPassword = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.MemberID == foreign.ID {
				local.R.Member = foreign
				if foreign.R == nil {
					foreign.R = &memberR{}
				}
				foreign.R.FallbackPassword = local
				break
			}
		}
	}

	return nil
}

// SetMember of the fallbackPassword to the related item.
// Sets o.R.Member to related.
// Adds o to related.R.FallbackPassword.
func (o *FallbackPassword) SetMember(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Member) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"fallback_passwords\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, []string{"member_id"}),
		strmangle.WhereClause("\"", "\"", 0, fallbackPasswordPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.MemberID = related.ID
	if o.R == nil {
		o.R = &fallbackPasswordR{
			Member: related,
		}
	} else {
		o.R.Member = related
	}

	if related.R == nil {
		related.R = &memberR{
			FallbackPassword: o,
		}
	} else {
		related.R.FallbackPassword = o
	}

	return nil
}

// FallbackPasswords retrieves all the records using an executor.
func FallbackPasswords(mods ...qm.QueryMod) fallbackPasswordQuery {
	mods = append(mods, qm.From("\"fallback_passwords\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"fallback_passwords\".*"})
	}

	return fallbackPasswordQuery{q}
}

// FindFallbackPassword retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFallbackPassword(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*FallbackPassword, error) {
	fallbackPasswordObj := &FallbackPassword{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"fallback_passwords\" where \"id\"=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, fallbackPasswordObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from fallback_passwords")
	}

	if err = fallbackPasswordObj.doAfterSelectHooks(ctx, exec); err != nil {
		return fallbackPasswordObj, err
	}

	return fallbackPasswordObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *FallbackPassword) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no fallback_passwords provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(fallbackPasswordColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	fallbackPasswordInsertCacheMut.RLock()
	cache, cached := fallbackPasswordInsertCache[key]
	fallbackPasswordInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			fallbackPasswordAllColumns,
			fallbackPasswordColumnsWithDefault,
			fallbackPasswordColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(fallbackPasswordType, fallbackPasswordMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(fallbackPasswordType, fallbackPasswordMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"fallback_passwords\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"fallback_passwords\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into fallback_passwords")
	}

	if !cached {
		fallbackPasswordInsertCacheMut.Lock()
		fallbackPasswordInsertCache[key] = cache
		fallbackPasswordInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the FallbackPassword.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *FallbackPassword) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	fallbackPasswordUpdateCacheMut.RLock()
	cache, cached := fallbackPasswordUpdateCache[key]
	fallbackPasswordUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			fallbackPasswordAllColumns,
			fallbackPasswordPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update fallback_passwords, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"fallback_passwords\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, fallbackPasswordPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(fallbackPasswordType, fallbackPasswordMapping, append(wl, fallbackPasswordPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update fallback_passwords row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for fallback_passwords")
	}

	if !cached {
		fallbackPasswordUpdateCacheMut.Lock()
		fallbackPasswordUpdateCache[key] = cache
		fallbackPasswordUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q fallbackPasswordQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for fallback_passwords")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for fallback_passwords")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FallbackPasswordSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fallbackPasswordPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"fallback_passwords\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, fallbackPasswordPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in fallbackPassword slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all fallbackPassword")
	}
	return rowsAff, nil
}

// Delete deletes a single FallbackPassword record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *FallbackPassword) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no FallbackPassword provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), fallbackPasswordPrimaryKeyMapping)
	sql := "DELETE FROM \"fallback_passwords\" WHERE \"id\"=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from fallback_passwords")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for fallback_passwords")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q fallbackPasswordQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no fallbackPasswordQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from fallback_passwords")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for fallback_passwords")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FallbackPasswordSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(fallbackPasswordBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fallbackPasswordPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"fallback_passwords\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, fallbackPasswordPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from fallbackPassword slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for fallback_passwords")
	}

	if len(fallbackPasswordAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *FallbackPassword) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindFallbackPassword(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FallbackPasswordSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := FallbackPasswordSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fallbackPasswordPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"fallback_passwords\".* FROM \"fallback_passwords\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, fallbackPasswordPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in FallbackPasswordSlice")
	}

	*o = slice

	return nil
}

// FallbackPasswordExists checks if the FallbackPassword row exists.
func FallbackPasswordExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"fallback_passwords\" where \"id\"=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if fallback_passwords exists")
	}

	return exists, nil
}

// Exists checks if the FallbackPassword row exists.
func (o *FallbackPassword) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return FallbackPasswordExists(ctx, exec, o.ID)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *FallbackPassword) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no fallback_passwords provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(fallbackPasswordColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	fallbackPasswordUpsertCacheMut.RLock()
	cache, cached := fallbackPasswordUpsertCache[key]
	fallbackPasswordUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			fallbackPasswordAllColumns,
			fallbackPasswordColumnsWithDefault,
			fallbackPasswordColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			fallbackPasswordAllColumns,
			fallbackPasswordPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert fallback_passwords, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(fallbackPasswordPrimaryKeyColumns))
			copy(conflict, fallbackPasswordPrimaryKeyColumns)
		}
		cache.query = buildUpsertQuerySQLite(dialect, "\"fallback_passwords\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(fallbackPasswordType, fallbackPasswordMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(fallbackPasswordType, fallbackPasswordMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert fallback_passwords")
	}

	if !cached {
		fallbackPasswordUpsertCacheMut.Lock()
		fallbackPasswordUpsertCache[key] = cache
		fallbackPasswordUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
