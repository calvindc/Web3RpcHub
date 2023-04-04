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

// FallbackResetToken is an object representing the database table.
type FallbackResetToken struct {
	ID          int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	HashedToken string    `boil:"hashed_token" json:"hashed_token" toml:"hashed_token" yaml:"hashed_token"`
	CreatedBy   int64     `boil:"created_by" json:"created_by" toml:"created_by" yaml:"created_by"`
	CreatedAt   time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	ForMember   int64     `boil:"for_member" json:"for_member" toml:"for_member" yaml:"for_member"`
	Active      bool      `boil:"active" json:"active" toml:"active" yaml:"active"`

	R *fallbackResetTokenR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L fallbackResetTokenL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var FallbackResetTokenColumns = struct {
	ID          string
	HashedToken string
	CreatedBy   string
	CreatedAt   string
	ForMember   string
	Active      string
}{
	ID:          "id",
	HashedToken: "hashed_token",
	CreatedBy:   "created_by",
	CreatedAt:   "created_at",
	ForMember:   "for_member",
	Active:      "active",
}

var FallbackResetTokenTableColumns = struct {
	ID          string
	HashedToken string
	CreatedBy   string
	CreatedAt   string
	ForMember   string
	Active      string
}{
	ID:          "fallback_reset_tokens.id",
	HashedToken: "fallback_reset_tokens.hashed_token",
	CreatedBy:   "fallback_reset_tokens.created_by",
	CreatedAt:   "fallback_reset_tokens.created_at",
	ForMember:   "fallback_reset_tokens.for_member",
	Active:      "fallback_reset_tokens.active",
}

// Generated where

var FallbackResetTokenWhere = struct {
	ID          whereHelperint64
	HashedToken whereHelperstring
	CreatedBy   whereHelperint64
	CreatedAt   whereHelpertime_Time
	ForMember   whereHelperint64
	Active      whereHelperbool
}{
	ID:          whereHelperint64{field: "\"fallback_reset_tokens\".\"id\""},
	HashedToken: whereHelperstring{field: "\"fallback_reset_tokens\".\"hashed_token\""},
	CreatedBy:   whereHelperint64{field: "\"fallback_reset_tokens\".\"created_by\""},
	CreatedAt:   whereHelpertime_Time{field: "\"fallback_reset_tokens\".\"created_at\""},
	ForMember:   whereHelperint64{field: "\"fallback_reset_tokens\".\"for_member\""},
	Active:      whereHelperbool{field: "\"fallback_reset_tokens\".\"active\""},
}

// FallbackResetTokenRels is where relationship names are stored.
var FallbackResetTokenRels = struct {
	ForMemberMember string
	CreatedByMember string
}{
	ForMemberMember: "ForMemberMember",
	CreatedByMember: "CreatedByMember",
}

// fallbackResetTokenR is where relationships are stored.
type fallbackResetTokenR struct {
	ForMemberMember *Member `boil:"ForMemberMember" json:"ForMemberMember" toml:"ForMemberMember" yaml:"ForMemberMember"`
	CreatedByMember *Member `boil:"CreatedByMember" json:"CreatedByMember" toml:"CreatedByMember" yaml:"CreatedByMember"`
}

// NewStruct creates a new relationship struct
func (*fallbackResetTokenR) NewStruct() *fallbackResetTokenR {
	return &fallbackResetTokenR{}
}

func (r *fallbackResetTokenR) GetForMemberMember() *Member {
	if r == nil {
		return nil
	}
	return r.ForMemberMember
}

func (r *fallbackResetTokenR) GetCreatedByMember() *Member {
	if r == nil {
		return nil
	}
	return r.CreatedByMember
}

// fallbackResetTokenL is where Load methods for each relationship are stored.
type fallbackResetTokenL struct{}

var (
	fallbackResetTokenAllColumns            = []string{"id", "hashed_token", "created_by", "created_at", "for_member", "active"}
	fallbackResetTokenColumnsWithoutDefault = []string{}
	fallbackResetTokenColumnsWithDefault    = []string{"id", "hashed_token", "created_by", "created_at", "for_member", "active"}
	fallbackResetTokenPrimaryKeyColumns     = []string{"id"}
	fallbackResetTokenGeneratedColumns      = []string{}
)

type (
	// FallbackResetTokenSlice is an alias for a slice of pointers to FallbackResetToken.
	// This should almost always be used instead of []FallbackResetToken.
	FallbackResetTokenSlice []*FallbackResetToken
	// FallbackResetTokenHook is the signature for custom FallbackResetToken hook methods
	FallbackResetTokenHook func(context.Context, boil.ContextExecutor, *FallbackResetToken) error

	fallbackResetTokenQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	fallbackResetTokenType                 = reflect.TypeOf(&FallbackResetToken{})
	fallbackResetTokenMapping              = queries.MakeStructMapping(fallbackResetTokenType)
	fallbackResetTokenPrimaryKeyMapping, _ = queries.BindMapping(fallbackResetTokenType, fallbackResetTokenMapping, fallbackResetTokenPrimaryKeyColumns)
	fallbackResetTokenInsertCacheMut       sync.RWMutex
	fallbackResetTokenInsertCache          = make(map[string]insertCache)
	fallbackResetTokenUpdateCacheMut       sync.RWMutex
	fallbackResetTokenUpdateCache          = make(map[string]updateCache)
	fallbackResetTokenUpsertCacheMut       sync.RWMutex
	fallbackResetTokenUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var fallbackResetTokenAfterSelectHooks []FallbackResetTokenHook

var fallbackResetTokenBeforeInsertHooks []FallbackResetTokenHook
var fallbackResetTokenAfterInsertHooks []FallbackResetTokenHook

var fallbackResetTokenBeforeUpdateHooks []FallbackResetTokenHook
var fallbackResetTokenAfterUpdateHooks []FallbackResetTokenHook

var fallbackResetTokenBeforeDeleteHooks []FallbackResetTokenHook
var fallbackResetTokenAfterDeleteHooks []FallbackResetTokenHook

var fallbackResetTokenBeforeUpsertHooks []FallbackResetTokenHook
var fallbackResetTokenAfterUpsertHooks []FallbackResetTokenHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *FallbackResetToken) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *FallbackResetToken) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *FallbackResetToken) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *FallbackResetToken) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *FallbackResetToken) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *FallbackResetToken) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *FallbackResetToken) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *FallbackResetToken) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *FallbackResetToken) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range fallbackResetTokenAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddFallbackResetTokenHook registers your hook function for all future operations.
func AddFallbackResetTokenHook(hookPoint boil.HookPoint, fallbackResetTokenHook FallbackResetTokenHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		fallbackResetTokenAfterSelectHooks = append(fallbackResetTokenAfterSelectHooks, fallbackResetTokenHook)
	case boil.BeforeInsertHook:
		fallbackResetTokenBeforeInsertHooks = append(fallbackResetTokenBeforeInsertHooks, fallbackResetTokenHook)
	case boil.AfterInsertHook:
		fallbackResetTokenAfterInsertHooks = append(fallbackResetTokenAfterInsertHooks, fallbackResetTokenHook)
	case boil.BeforeUpdateHook:
		fallbackResetTokenBeforeUpdateHooks = append(fallbackResetTokenBeforeUpdateHooks, fallbackResetTokenHook)
	case boil.AfterUpdateHook:
		fallbackResetTokenAfterUpdateHooks = append(fallbackResetTokenAfterUpdateHooks, fallbackResetTokenHook)
	case boil.BeforeDeleteHook:
		fallbackResetTokenBeforeDeleteHooks = append(fallbackResetTokenBeforeDeleteHooks, fallbackResetTokenHook)
	case boil.AfterDeleteHook:
		fallbackResetTokenAfterDeleteHooks = append(fallbackResetTokenAfterDeleteHooks, fallbackResetTokenHook)
	case boil.BeforeUpsertHook:
		fallbackResetTokenBeforeUpsertHooks = append(fallbackResetTokenBeforeUpsertHooks, fallbackResetTokenHook)
	case boil.AfterUpsertHook:
		fallbackResetTokenAfterUpsertHooks = append(fallbackResetTokenAfterUpsertHooks, fallbackResetTokenHook)
	}
}

// One returns a single fallbackResetToken record from the query.
func (q fallbackResetTokenQuery) One(ctx context.Context, exec boil.ContextExecutor) (*FallbackResetToken, error) {
	o := &FallbackResetToken{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for fallback_reset_tokens")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all FallbackResetToken records from the query.
func (q fallbackResetTokenQuery) All(ctx context.Context, exec boil.ContextExecutor) (FallbackResetTokenSlice, error) {
	var o []*FallbackResetToken

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to FallbackResetToken slice")
	}

	if len(fallbackResetTokenAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all FallbackResetToken records in the query.
func (q fallbackResetTokenQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count fallback_reset_tokens rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q fallbackResetTokenQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if fallback_reset_tokens exists")
	}

	return count > 0, nil
}

// ForMemberMember pointed to by the foreign key.
func (o *FallbackResetToken) ForMemberMember(mods ...qm.QueryMod) memberQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ForMember),
	}

	queryMods = append(queryMods, mods...)

	return Members(queryMods...)
}

// CreatedByMember pointed to by the foreign key.
func (o *FallbackResetToken) CreatedByMember(mods ...qm.QueryMod) memberQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.CreatedBy),
	}

	queryMods = append(queryMods, mods...)

	return Members(queryMods...)
}

// LoadForMemberMember allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (fallbackResetTokenL) LoadForMemberMember(ctx context.Context, e boil.ContextExecutor, singular bool, maybeFallbackResetToken interface{}, mods queries.Applicator) error {
	var slice []*FallbackResetToken
	var object *FallbackResetToken

	if singular {
		var ok bool
		object, ok = maybeFallbackResetToken.(*FallbackResetToken)
		if !ok {
			object = new(FallbackResetToken)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeFallbackResetToken)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeFallbackResetToken))
			}
		}
	} else {
		s, ok := maybeFallbackResetToken.(*[]*FallbackResetToken)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeFallbackResetToken)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeFallbackResetToken))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &fallbackResetTokenR{}
		}
		args = append(args, object.ForMember)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &fallbackResetTokenR{}
			}

			for _, a := range args {
				if a == obj.ForMember {
					continue Outer
				}
			}

			args = append(args, obj.ForMember)

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
		object.R.ForMemberMember = foreign
		if foreign.R == nil {
			foreign.R = &memberR{}
		}
		foreign.R.ForMemberFallbackResetTokens = append(foreign.R.ForMemberFallbackResetTokens, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ForMember == foreign.ID {
				local.R.ForMemberMember = foreign
				if foreign.R == nil {
					foreign.R = &memberR{}
				}
				foreign.R.ForMemberFallbackResetTokens = append(foreign.R.ForMemberFallbackResetTokens, local)
				break
			}
		}
	}

	return nil
}

// LoadCreatedByMember allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (fallbackResetTokenL) LoadCreatedByMember(ctx context.Context, e boil.ContextExecutor, singular bool, maybeFallbackResetToken interface{}, mods queries.Applicator) error {
	var slice []*FallbackResetToken
	var object *FallbackResetToken

	if singular {
		var ok bool
		object, ok = maybeFallbackResetToken.(*FallbackResetToken)
		if !ok {
			object = new(FallbackResetToken)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeFallbackResetToken)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeFallbackResetToken))
			}
		}
	} else {
		s, ok := maybeFallbackResetToken.(*[]*FallbackResetToken)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeFallbackResetToken)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeFallbackResetToken))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &fallbackResetTokenR{}
		}
		args = append(args, object.CreatedBy)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &fallbackResetTokenR{}
			}

			for _, a := range args {
				if a == obj.CreatedBy {
					continue Outer
				}
			}

			args = append(args, obj.CreatedBy)

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
		object.R.CreatedByMember = foreign
		if foreign.R == nil {
			foreign.R = &memberR{}
		}
		foreign.R.CreatedByFallbackResetTokens = append(foreign.R.CreatedByFallbackResetTokens, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CreatedBy == foreign.ID {
				local.R.CreatedByMember = foreign
				if foreign.R == nil {
					foreign.R = &memberR{}
				}
				foreign.R.CreatedByFallbackResetTokens = append(foreign.R.CreatedByFallbackResetTokens, local)
				break
			}
		}
	}

	return nil
}

// SetForMemberMember of the fallbackResetToken to the related item.
// Sets o.R.ForMemberMember to related.
// Adds o to related.R.ForMemberFallbackResetTokens.
func (o *FallbackResetToken) SetForMemberMember(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Member) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"fallback_reset_tokens\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, []string{"for_member"}),
		strmangle.WhereClause("\"", "\"", 0, fallbackResetTokenPrimaryKeyColumns),
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

	o.ForMember = related.ID
	if o.R == nil {
		o.R = &fallbackResetTokenR{
			ForMemberMember: related,
		}
	} else {
		o.R.ForMemberMember = related
	}

	if related.R == nil {
		related.R = &memberR{
			ForMemberFallbackResetTokens: FallbackResetTokenSlice{o},
		}
	} else {
		related.R.ForMemberFallbackResetTokens = append(related.R.ForMemberFallbackResetTokens, o)
	}

	return nil
}

// SetCreatedByMember of the fallbackResetToken to the related item.
// Sets o.R.CreatedByMember to related.
// Adds o to related.R.CreatedByFallbackResetTokens.
func (o *FallbackResetToken) SetCreatedByMember(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Member) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"fallback_reset_tokens\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, []string{"created_by"}),
		strmangle.WhereClause("\"", "\"", 0, fallbackResetTokenPrimaryKeyColumns),
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

	o.CreatedBy = related.ID
	if o.R == nil {
		o.R = &fallbackResetTokenR{
			CreatedByMember: related,
		}
	} else {
		o.R.CreatedByMember = related
	}

	if related.R == nil {
		related.R = &memberR{
			CreatedByFallbackResetTokens: FallbackResetTokenSlice{o},
		}
	} else {
		related.R.CreatedByFallbackResetTokens = append(related.R.CreatedByFallbackResetTokens, o)
	}

	return nil
}

// FallbackResetTokens retrieves all the records using an executor.
func FallbackResetTokens(mods ...qm.QueryMod) fallbackResetTokenQuery {
	mods = append(mods, qm.From("\"fallback_reset_tokens\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"fallback_reset_tokens\".*"})
	}

	return fallbackResetTokenQuery{q}
}

// FindFallbackResetToken retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFallbackResetToken(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*FallbackResetToken, error) {
	fallbackResetTokenObj := &FallbackResetToken{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"fallback_reset_tokens\" where \"id\"=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, fallbackResetTokenObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from fallback_reset_tokens")
	}

	if err = fallbackResetTokenObj.doAfterSelectHooks(ctx, exec); err != nil {
		return fallbackResetTokenObj, err
	}

	return fallbackResetTokenObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *FallbackResetToken) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no fallback_reset_tokens provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(fallbackResetTokenColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	fallbackResetTokenInsertCacheMut.RLock()
	cache, cached := fallbackResetTokenInsertCache[key]
	fallbackResetTokenInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			fallbackResetTokenAllColumns,
			fallbackResetTokenColumnsWithDefault,
			fallbackResetTokenColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(fallbackResetTokenType, fallbackResetTokenMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(fallbackResetTokenType, fallbackResetTokenMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"fallback_reset_tokens\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"fallback_reset_tokens\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into fallback_reset_tokens")
	}

	if !cached {
		fallbackResetTokenInsertCacheMut.Lock()
		fallbackResetTokenInsertCache[key] = cache
		fallbackResetTokenInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the FallbackResetToken.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *FallbackResetToken) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	fallbackResetTokenUpdateCacheMut.RLock()
	cache, cached := fallbackResetTokenUpdateCache[key]
	fallbackResetTokenUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			fallbackResetTokenAllColumns,
			fallbackResetTokenPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update fallback_reset_tokens, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"fallback_reset_tokens\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, fallbackResetTokenPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(fallbackResetTokenType, fallbackResetTokenMapping, append(wl, fallbackResetTokenPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update fallback_reset_tokens row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for fallback_reset_tokens")
	}

	if !cached {
		fallbackResetTokenUpdateCacheMut.Lock()
		fallbackResetTokenUpdateCache[key] = cache
		fallbackResetTokenUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q fallbackResetTokenQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for fallback_reset_tokens")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for fallback_reset_tokens")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FallbackResetTokenSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fallbackResetTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"fallback_reset_tokens\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, fallbackResetTokenPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in fallbackResetToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all fallbackResetToken")
	}
	return rowsAff, nil
}

// Delete deletes a single FallbackResetToken record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *FallbackResetToken) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no FallbackResetToken provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), fallbackResetTokenPrimaryKeyMapping)
	sql := "DELETE FROM \"fallback_reset_tokens\" WHERE \"id\"=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from fallback_reset_tokens")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for fallback_reset_tokens")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q fallbackResetTokenQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no fallbackResetTokenQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from fallback_reset_tokens")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for fallback_reset_tokens")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FallbackResetTokenSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(fallbackResetTokenBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fallbackResetTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"fallback_reset_tokens\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, fallbackResetTokenPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from fallbackResetToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for fallback_reset_tokens")
	}

	if len(fallbackResetTokenAfterDeleteHooks) != 0 {
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
func (o *FallbackResetToken) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindFallbackResetToken(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FallbackResetTokenSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := FallbackResetTokenSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fallbackResetTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"fallback_reset_tokens\".* FROM \"fallback_reset_tokens\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, fallbackResetTokenPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in FallbackResetTokenSlice")
	}

	*o = slice

	return nil
}

// FallbackResetTokenExists checks if the FallbackResetToken row exists.
func FallbackResetTokenExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"fallback_reset_tokens\" where \"id\"=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if fallback_reset_tokens exists")
	}

	return exists, nil
}

// Exists checks if the FallbackResetToken row exists.
func (o *FallbackResetToken) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return FallbackResetTokenExists(ctx, exec, o.ID)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *FallbackResetToken) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no fallback_reset_tokens provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(fallbackResetTokenColumnsWithDefault, o)

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

	fallbackResetTokenUpsertCacheMut.RLock()
	cache, cached := fallbackResetTokenUpsertCache[key]
	fallbackResetTokenUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			fallbackResetTokenAllColumns,
			fallbackResetTokenColumnsWithDefault,
			fallbackResetTokenColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			fallbackResetTokenAllColumns,
			fallbackResetTokenPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert fallback_reset_tokens, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(fallbackResetTokenPrimaryKeyColumns))
			copy(conflict, fallbackResetTokenPrimaryKeyColumns)
		}
		cache.query = buildUpsertQuerySQLite(dialect, "\"fallback_reset_tokens\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(fallbackResetTokenType, fallbackResetTokenMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(fallbackResetTokenType, fallbackResetTokenMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert fallback_reset_tokens")
	}

	if !cached {
		fallbackResetTokenUpsertCacheMut.Lock()
		fallbackResetTokenUpsertCache[key] = cache
		fallbackResetTokenUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
