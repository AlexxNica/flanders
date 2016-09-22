package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
)

// Alia is an object representing the database table.
type Alia struct {
	ID    int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	Key   string `boil:"key" json:"key" toml:"key" yaml:"key"`
	Value string `boil:"value" json:"value" toml:"value" yaml:"value"`

	R *aliaR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L aliaL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// aliaR is where relationships are stored.
type aliaR struct {
}

// aliaL is where Load methods for each relationship are stored.
type aliaL struct{}

var (
	aliaColumns               = []string{"id", "key", "value"}
	aliaColumnsWithoutDefault = []string{"key", "value"}
	aliaColumnsWithDefault    = []string{"id"}
	aliaPrimaryKeyColumns     = []string{"id"}
)

type (
	// AliaSlice is an alias for a slice of pointers to Alia.
	// This should generally be used opposed to []Alia.
	AliaSlice []*Alia
	// AliaHook is the signature for custom Alia hook methods
	AliaHook func(boil.Executor, *Alia) error

	aliaQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	aliaType                 = reflect.TypeOf(&Alia{})
	aliaMapping              = queries.MakeStructMapping(aliaType)
	aliaPrimaryKeyMapping, _ = queries.BindMapping(aliaType, aliaMapping, aliaPrimaryKeyColumns)
	aliaInsertCacheMut       sync.RWMutex
	aliaInsertCache          = make(map[string]insertCache)
	aliaUpdateCacheMut       sync.RWMutex
	aliaUpdateCache          = make(map[string]updateCache)
	aliaUpsertCacheMut       sync.RWMutex
	aliaUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

var aliaBeforeInsertHooks []AliaHook
var aliaBeforeUpdateHooks []AliaHook
var aliaBeforeDeleteHooks []AliaHook
var aliaBeforeUpsertHooks []AliaHook

var aliaAfterInsertHooks []AliaHook
var aliaAfterSelectHooks []AliaHook
var aliaAfterUpdateHooks []AliaHook
var aliaAfterDeleteHooks []AliaHook
var aliaAfterUpsertHooks []AliaHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Alia) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Alia) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Alia) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Alia) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Alia) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Alia) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Alia) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Alia) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Alia) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range aliaAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddAliaHook registers your hook function for all future operations.
func AddAliaHook(hookPoint boil.HookPoint, aliaHook AliaHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		aliaBeforeInsertHooks = append(aliaBeforeInsertHooks, aliaHook)
	case boil.BeforeUpdateHook:
		aliaBeforeUpdateHooks = append(aliaBeforeUpdateHooks, aliaHook)
	case boil.BeforeDeleteHook:
		aliaBeforeDeleteHooks = append(aliaBeforeDeleteHooks, aliaHook)
	case boil.BeforeUpsertHook:
		aliaBeforeUpsertHooks = append(aliaBeforeUpsertHooks, aliaHook)
	case boil.AfterInsertHook:
		aliaAfterInsertHooks = append(aliaAfterInsertHooks, aliaHook)
	case boil.AfterSelectHook:
		aliaAfterSelectHooks = append(aliaAfterSelectHooks, aliaHook)
	case boil.AfterUpdateHook:
		aliaAfterUpdateHooks = append(aliaAfterUpdateHooks, aliaHook)
	case boil.AfterDeleteHook:
		aliaAfterDeleteHooks = append(aliaAfterDeleteHooks, aliaHook)
	case boil.AfterUpsertHook:
		aliaAfterUpsertHooks = append(aliaAfterUpsertHooks, aliaHook)
	}
}

// OneP returns a single alia record from the query, and panics on error.
func (q aliaQuery) OneP() *Alia {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single alia record from the query.
func (q aliaQuery) One() (*Alia, error) {
	o := &Alia{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for alias")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Alia records from the query, and panics on error.
func (q aliaQuery) AllP() AliaSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Alia records from the query.
func (q aliaQuery) All() (AliaSlice, error) {
	var o AliaSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Alia slice")
	}

	if len(aliaAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Alia records in the query, and panics on error.
func (q aliaQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Alia records in the query.
func (q aliaQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count alias rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q aliaQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q aliaQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if alias exists")
	}

	return count > 0, nil
}

















// AliasesG retrieves all records.
func AliasesG(mods ...qm.QueryMod) aliaQuery {
	return Aliases(boil.GetDB(), mods...)
}

// Aliases retrieves all the records using an executor.
func Aliases(exec boil.Executor, mods ...qm.QueryMod) aliaQuery {
	mods = append(mods, qm.From("`alias`"))
	return aliaQuery{NewQuery(exec, mods...)}
}

// FindAliaG retrieves a single record by ID.
func FindAliaG(id int64, selectCols ...string) (*Alia, error) {
	return FindAlia(boil.GetDB(), id, selectCols...)
}

// FindAliaGP retrieves a single record by ID, and panics on error.
func FindAliaGP(id int64, selectCols ...string) *Alia {
	retobj, err := FindAlia(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindAlia retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAlia(exec boil.Executor, id int64, selectCols ...string) (*Alia, error) {
	aliaObj := &Alia{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `alias` where `id`=?", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(aliaObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from alias")
	}

	return aliaObj, nil
}

// FindAliaP retrieves a single record by ID with an executor, and panics on error.
func FindAliaP(exec boil.Executor, id int64, selectCols ...string) *Alia {
	retobj, err := FindAlia(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Alia) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Alia) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Alia) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Alia) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no alias provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(aliaColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	aliaInsertCacheMut.RLock()
	cache, cached := aliaInsertCache[key]
	aliaInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			aliaColumns,
			aliaColumnsWithDefault,
			aliaColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(aliaType, aliaMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(aliaType, aliaMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO `alias` (`%s`) VALUES (%s)", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `alias` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, aliaPrimaryKeyColumns))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.Exec(cache.query, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into alias")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int64(lastID)
	identifierCols = []interface{}{lastID}

	if lastID == 0 || len(cache.retMapping) != 1 || cache.retMapping[0] == aliaMapping["ID"] {
		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, cache.retQuery)
			fmt.Fprintln(boil.DebugWriter, identifierCols...)
		}

		err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
		if err != nil {
			return errors.Wrap(err, "models: unable to populate default values for alias")
		}
	}

CacheNoHooks:
	if !cached {
		aliaInsertCacheMut.Lock()
		aliaInsertCache[key] = cache
		aliaInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Alia record. See Update for
// whitelist behavior description.
func (o *Alia) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Alia record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Alia) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Alia, and panics on error.
// See Update for whitelist behavior description.
func (o *Alia) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Alia.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Alia) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	aliaUpdateCacheMut.RLock()
	cache, cached := aliaUpdateCache[key]
	aliaUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(aliaColumns, aliaPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update alias, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `alias` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, aliaPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(aliaType, aliaMapping, append(wl, aliaPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update alias row")
	}

	if !cached {
		aliaUpdateCacheMut.Lock()
		aliaUpdateCache[key] = cache
		aliaUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q aliaQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q aliaQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for alias")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o AliaSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o AliaSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o AliaSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AliaSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), aliaPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE `alias` SET %s WHERE (`id`) IN (%s)",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(aliaPrimaryKeyColumns), len(colNames)+1, len(aliaPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in alia slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Alia) UpsertG(updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Alia) UpsertGP(updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Alia) UpsertP(exec boil.Executor, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Alia) Upsert(exec boil.Executor, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no alias provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(aliaColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	aliaUpsertCacheMut.RLock()
	cache, cached := aliaUpsertCache[key]
	aliaUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			aliaColumns,
			aliaColumnsWithDefault,
			aliaColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			aliaColumns,
			aliaPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert alias, could not build update column list")
		}

		cache.query = queries.BuildUpsertQueryMySQL(dialect, "alias", update, whitelist)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `alias` WHERE `id`=?",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
		)

		cache.valueMapping, err = queries.BindMapping(aliaType, aliaMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(aliaType, aliaMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	values := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	result, err := exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for alias")
	}

	var lastID int64
	var identifierCols []interface{}
	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int64(lastID)
	identifierCols = []interface{}{lastID}

	if lastID == 0 || len(cache.retMapping) != 1 || cache.retMapping[0] == aliaMapping["ID"] {
		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, cache.retQuery)
			fmt.Fprintln(boil.DebugWriter, identifierCols...)
		}

		err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(returns...)
		if err != nil {
			return errors.Wrap(err, "models: unable to populate default values for alias")
		}
	}

CacheNoHooks:
	if !cached {
		aliaUpsertCacheMut.Lock()
		aliaUpsertCache[key] = cache
		aliaUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Alia record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Alia) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Alia record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Alia) DeleteG() error {
	if o == nil {
		return errors.New("models: no Alia provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Alia record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Alia) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Alia record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Alia) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Alia provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), aliaPrimaryKeyMapping)
	sql := "DELETE FROM `alias` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from alias")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q aliaQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q aliaQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no aliaQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from alias")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o AliaSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o AliaSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Alia slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o AliaSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AliaSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Alia slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(aliaBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), aliaPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM `alias` WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, aliaPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(aliaPrimaryKeyColumns), 1, len(aliaPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from alia slice")
	}

	if len(aliaAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Alia) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Alia) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Alia) ReloadG() error {
	if o == nil {
		return errors.New("models: no Alia provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Alia) Reload(exec boil.Executor) error {
	ret, err := FindAlia(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *AliaSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *AliaSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AliaSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty AliaSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AliaSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	aliases := AliaSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), aliaPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT `alias`.* FROM `alias` WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, aliaPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(aliaPrimaryKeyColumns), 1, len(aliaPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&aliases)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AliaSlice")
	}

	*o = aliases

	return nil
}

// AliaExists checks if the Alia row exists.
func AliaExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from `alias` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if alias exists")
	}

	return exists, nil
}

// AliaExistsG checks if the Alia row exists.
func AliaExistsG(id int64) (bool, error) {
	return AliaExists(boil.GetDB(), id)
}

// AliaExistsGP checks if the Alia row exists. Panics on error.
func AliaExistsGP(id int64) bool {
	e, err := AliaExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// AliaExistsP checks if the Alia row exists. Panics on error.
func AliaExistsP(exec boil.Executor, id int64) bool {
	e, err := AliaExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}


