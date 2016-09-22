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
	"gopkg.in/nullbio/null.v5"
)

// Message is an object representing the database table.
type Message struct {
	ID              int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	Date            time.Time `boil:"date" json:"date" toml:"date" yaml:"date"`
	MicroTS         int64     `boil:"micro_ts" json:"micro_ts" toml:"micro_ts" yaml:"micro_ts"`
	Method          string    `boil:"method" json:"method" toml:"method" yaml:"method"`
	ReplyReason     string    `boil:"reply_reason" json:"reply_reason" toml:"reply_reason" yaml:"reply_reason"`
	Ruri            string    `boil:"ruri" json:"ruri" toml:"ruri" yaml:"ruri"`
	RuriUser        string    `boil:"ruri_user" json:"ruri_user" toml:"ruri_user" yaml:"ruri_user"`
	RuriDomain      string    `boil:"ruri_domain" json:"ruri_domain" toml:"ruri_domain" yaml:"ruri_domain"`
	FromUser        string    `boil:"from_user" json:"from_user" toml:"from_user" yaml:"from_user"`
	FromDomain      string    `boil:"from_domain" json:"from_domain" toml:"from_domain" yaml:"from_domain"`
	FromTag         string    `boil:"from_tag" json:"from_tag" toml:"from_tag" yaml:"from_tag"`
	ToUser          string    `boil:"to_user" json:"to_user" toml:"to_user" yaml:"to_user"`
	ToDomain        string    `boil:"to_domain" json:"to_domain" toml:"to_domain" yaml:"to_domain"`
	ToTag           string    `boil:"to_tag" json:"to_tag" toml:"to_tag" yaml:"to_tag"`
	PidUser         string    `boil:"pid_user" json:"pid_user" toml:"pid_user" yaml:"pid_user"`
	ContactUser     string    `boil:"contact_user" json:"contact_user" toml:"contact_user" yaml:"contact_user"`
	AuthUser        string    `boil:"auth_user" json:"auth_user" toml:"auth_user" yaml:"auth_user"`
	Callid          string    `boil:"callid" json:"callid" toml:"callid" yaml:"callid"`
	CallidAleg      string    `boil:"callid_aleg" json:"callid_aleg" toml:"callid_aleg" yaml:"callid_aleg"`
	Via1            string    `boil:"via_1" json:"via_1" toml:"via_1" yaml:"via_1"`
	Via1Branch      string    `boil:"via_1_branch" json:"via_1_branch" toml:"via_1_branch" yaml:"via_1_branch"`
	Cseq            string    `boil:"cseq" json:"cseq" toml:"cseq" yaml:"cseq"`
	Diversion       string    `boil:"diversion" json:"diversion" toml:"diversion" yaml:"diversion"`
	Reason          string    `boil:"reason" json:"reason" toml:"reason" yaml:"reason"`
	ContentType     string    `boil:"content_type" json:"content_type" toml:"content_type" yaml:"content_type"`
	Auth            string    `boil:"auth" json:"auth" toml:"auth" yaml:"auth"`
	UserAgent       string    `boil:"user_agent" json:"user_agent" toml:"user_agent" yaml:"user_agent"`
	SourceIp        string    `boil:"source_ip" json:"source_ip" toml:"source_ip" yaml:"source_ip"`
	SourcePort      int       `boil:"source_port" json:"source_port" toml:"source_port" yaml:"source_port"`
	DestinationIp   string    `boil:"destination_ip" json:"destination_ip" toml:"destination_ip" yaml:"destination_ip"`
	DestinationPort int       `boil:"destination_port" json:"destination_port" toml:"destination_port" yaml:"destination_port"`
	ContactIp       string    `boil:"contact_ip" json:"contact_ip" toml:"contact_ip" yaml:"contact_ip"`
	ContactPort     int       `boil:"contact_port" json:"contact_port" toml:"contact_port" yaml:"contact_port"`
	OriginatorIp    string    `boil:"originator_ip" json:"originator_ip" toml:"originator_ip" yaml:"originator_ip"`
	OriginatorPort  int       `boil:"originator_port" json:"originator_port" toml:"originator_port" yaml:"originator_port"`
	CorrelationID   string    `boil:"correlation_id" json:"correlation_id" toml:"correlation_id" yaml:"correlation_id"`
	CustomField1    string    `boil:"custom_field1" json:"custom_field1" toml:"custom_field1" yaml:"custom_field1"`
	CustomField2    string    `boil:"custom_field2" json:"custom_field2" toml:"custom_field2" yaml:"custom_field2"`
	CustomField3    string    `boil:"custom_field3" json:"custom_field3" toml:"custom_field3" yaml:"custom_field3"`
	Proto           int       `boil:"proto" json:"proto" toml:"proto" yaml:"proto"`
	Family          null.Int  `boil:"family" json:"family,omitempty" toml:"family" yaml:"family,omitempty"`
	RTPStat         string    `boil:"rtp_stat" json:"rtp_stat" toml:"rtp_stat" yaml:"rtp_stat"`
	Type            int       `boil:"type" json:"type" toml:"type" yaml:"type"`
	Node            string    `boil:"node" json:"node" toml:"node" yaml:"node"`
	MSG             string    `boil:"msg" json:"msg" toml:"msg" yaml:"msg"`

	R *messageR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L messageL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// messageR is where relationships are stored.
type messageR struct {
}

// messageL is where Load methods for each relationship are stored.
type messageL struct{}

var (
	messageColumns               = []string{"id", "date", "micro_ts", "method", "reply_reason", "ruri", "ruri_user", "ruri_domain", "from_user", "from_domain", "from_tag", "to_user", "to_domain", "to_tag", "pid_user", "contact_user", "auth_user", "callid", "callid_aleg", "via_1", "via_1_branch", "cseq", "diversion", "reason", "content_type", "auth", "user_agent", "source_ip", "source_port", "destination_ip", "destination_port", "contact_ip", "contact_port", "originator_ip", "originator_port", "correlation_id", "custom_field1", "custom_field2", "custom_field3", "proto", "family", "rtp_stat", "type", "node", "msg"}
	messageColumnsWithoutDefault = []string{"method", "reply_reason", "ruri", "ruri_user", "ruri_domain", "from_user", "from_domain", "from_tag", "to_user", "to_domain", "to_tag", "pid_user", "contact_user", "auth_user", "callid", "callid_aleg", "via_1", "via_1_branch", "cseq", "diversion", "reason", "content_type", "auth", "user_agent", "source_ip", "destination_ip", "contact_ip", "originator_ip", "correlation_id", "custom_field1", "custom_field2", "custom_field3", "family", "rtp_stat", "node", "msg"}
	messageColumnsWithDefault    = []string{"id", "date", "micro_ts", "source_port", "destination_port", "contact_port", "originator_port", "proto", "type"}
	messagePrimaryKeyColumns     = []string{"id", "date"}
)

type (
	// MessageSlice is an alias for a slice of pointers to Message.
	// This should generally be used opposed to []Message.
	MessageSlice []*Message
	// MessageHook is the signature for custom Message hook methods
	MessageHook func(boil.Executor, *Message) error

	messageQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	messageType                 = reflect.TypeOf(&Message{})
	messageMapping              = queries.MakeStructMapping(messageType)
	messagePrimaryKeyMapping, _ = queries.BindMapping(messageType, messageMapping, messagePrimaryKeyColumns)
	messageInsertCacheMut       sync.RWMutex
	messageInsertCache          = make(map[string]insertCache)
	messageUpdateCacheMut       sync.RWMutex
	messageUpdateCache          = make(map[string]updateCache)
	messageUpsertCacheMut       sync.RWMutex
	messageUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

var messageBeforeInsertHooks []MessageHook
var messageBeforeUpdateHooks []MessageHook
var messageBeforeDeleteHooks []MessageHook
var messageBeforeUpsertHooks []MessageHook

var messageAfterInsertHooks []MessageHook
var messageAfterSelectHooks []MessageHook
var messageAfterUpdateHooks []MessageHook
var messageAfterDeleteHooks []MessageHook
var messageAfterUpsertHooks []MessageHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Message) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range messageBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Message) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range messageBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Message) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range messageBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Message) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range messageBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Message) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range messageAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Message) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range messageAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Message) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range messageAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Message) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range messageAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Message) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range messageAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMessageHook registers your hook function for all future operations.
func AddMessageHook(hookPoint boil.HookPoint, messageHook MessageHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		messageBeforeInsertHooks = append(messageBeforeInsertHooks, messageHook)
	case boil.BeforeUpdateHook:
		messageBeforeUpdateHooks = append(messageBeforeUpdateHooks, messageHook)
	case boil.BeforeDeleteHook:
		messageBeforeDeleteHooks = append(messageBeforeDeleteHooks, messageHook)
	case boil.BeforeUpsertHook:
		messageBeforeUpsertHooks = append(messageBeforeUpsertHooks, messageHook)
	case boil.AfterInsertHook:
		messageAfterInsertHooks = append(messageAfterInsertHooks, messageHook)
	case boil.AfterSelectHook:
		messageAfterSelectHooks = append(messageAfterSelectHooks, messageHook)
	case boil.AfterUpdateHook:
		messageAfterUpdateHooks = append(messageAfterUpdateHooks, messageHook)
	case boil.AfterDeleteHook:
		messageAfterDeleteHooks = append(messageAfterDeleteHooks, messageHook)
	case boil.AfterUpsertHook:
		messageAfterUpsertHooks = append(messageAfterUpsertHooks, messageHook)
	}
}

// OneP returns a single message record from the query, and panics on error.
func (q messageQuery) OneP() *Message {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single message record from the query.
func (q messageQuery) One() (*Message, error) {
	o := &Message{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for message")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Message records from the query, and panics on error.
func (q messageQuery) AllP() MessageSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Message records from the query.
func (q messageQuery) All() (MessageSlice, error) {
	var o MessageSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Message slice")
	}

	if len(messageAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Message records in the query, and panics on error.
func (q messageQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Message records in the query.
func (q messageQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count message rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q messageQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q messageQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if message exists")
	}

	return count > 0, nil
}

















// MessagesG retrieves all records.
func MessagesG(mods ...qm.QueryMod) messageQuery {
	return Messages(boil.GetDB(), mods...)
}

// Messages retrieves all the records using an executor.
func Messages(exec boil.Executor, mods ...qm.QueryMod) messageQuery {
	mods = append(mods, qm.From("`message`"))
	return messageQuery{NewQuery(exec, mods...)}
}

// FindMessageG retrieves a single record by ID.
func FindMessageG(id int64, date time.Time, selectCols ...string) (*Message, error) {
	return FindMessage(boil.GetDB(), id, date, selectCols...)
}

// FindMessageGP retrieves a single record by ID, and panics on error.
func FindMessageGP(id int64, date time.Time, selectCols ...string) *Message {
	retobj, err := FindMessage(boil.GetDB(), id, date, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindMessage retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMessage(exec boil.Executor, id int64, date time.Time, selectCols ...string) (*Message, error) {
	messageObj := &Message{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `message` where `id`=? AND `date`=?", sel,
	)

	q := queries.Raw(exec, query, id, date)

	err := q.Bind(messageObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from message")
	}

	return messageObj, nil
}

// FindMessageP retrieves a single record by ID with an executor, and panics on error.
func FindMessageP(exec boil.Executor, id int64, date time.Time, selectCols ...string) *Message {
	retobj, err := FindMessage(exec, id, date, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Message) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Message) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Message) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Message) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no message provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(messageColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	messageInsertCacheMut.RLock()
	cache, cached := messageInsertCache[key]
	messageInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			messageColumns,
			messageColumnsWithDefault,
			messageColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(messageType, messageMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(messageType, messageMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO `message` (`%s`) VALUES (%s)", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `message` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, messagePrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into message")
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

	identifierCols = []interface{}{
		o.ID,
		o.Date,
	}

	if lastID == 0 || len(cache.retMapping) != 1 || cache.retMapping[0] == messageMapping["ID"] {
		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, cache.retQuery)
			fmt.Fprintln(boil.DebugWriter, identifierCols...)
		}

		err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
		if err != nil {
			return errors.Wrap(err, "models: unable to populate default values for message")
		}
	}

CacheNoHooks:
	if !cached {
		messageInsertCacheMut.Lock()
		messageInsertCache[key] = cache
		messageInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Message record. See Update for
// whitelist behavior description.
func (o *Message) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Message record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Message) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Message, and panics on error.
// See Update for whitelist behavior description.
func (o *Message) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Message.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Message) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	messageUpdateCacheMut.RLock()
	cache, cached := messageUpdateCache[key]
	messageUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(messageColumns, messagePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update message, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `message` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, messagePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(messageType, messageMapping, append(wl, messagePrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update message row")
	}

	if !cached {
		messageUpdateCacheMut.Lock()
		messageUpdateCache[key] = cache
		messageUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q messageQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q messageQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for message")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o MessageSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o MessageSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o MessageSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MessageSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), messagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE `message` SET %s WHERE (`id`,`date`) IN (%s)",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(messagePrimaryKeyColumns), len(colNames)+1, len(messagePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in message slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Message) UpsertG(updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Message) UpsertGP(updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Message) UpsertP(exec boil.Executor, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Message) Upsert(exec boil.Executor, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no message provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(messageColumnsWithDefault, o)

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

	messageUpsertCacheMut.RLock()
	cache, cached := messageUpsertCache[key]
	messageUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			messageColumns,
			messageColumnsWithDefault,
			messageColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			messageColumns,
			messagePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert message, could not build update column list")
		}

		cache.query = queries.BuildUpsertQueryMySQL(dialect, "message", update, whitelist)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `message` WHERE `id`=? AND `date`=?",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
		)

		cache.valueMapping, err = queries.BindMapping(messageType, messageMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(messageType, messageMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for message")
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

	identifierCols = []interface{}{
		o.ID,
		o.Date,
	}

	if lastID == 0 || len(cache.retMapping) != 1 || cache.retMapping[0] == messageMapping["ID"] {
		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, cache.retQuery)
			fmt.Fprintln(boil.DebugWriter, identifierCols...)
		}

		err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(returns...)
		if err != nil {
			return errors.Wrap(err, "models: unable to populate default values for message")
		}
	}

CacheNoHooks:
	if !cached {
		messageUpsertCacheMut.Lock()
		messageUpsertCache[key] = cache
		messageUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Message record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Message) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Message record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Message) DeleteG() error {
	if o == nil {
		return errors.New("models: no Message provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Message record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Message) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Message record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Message) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Message provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), messagePrimaryKeyMapping)
	sql := "DELETE FROM `message` WHERE `id`=? AND `date`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from message")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q messageQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q messageQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no messageQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from message")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o MessageSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o MessageSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Message slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o MessageSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MessageSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Message slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(messageBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), messagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM `message` WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, messagePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(messagePrimaryKeyColumns), 1, len(messagePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from message slice")
	}

	if len(messageAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Message) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Message) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Message) ReloadG() error {
	if o == nil {
		return errors.New("models: no Message provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Message) Reload(exec boil.Executor) error {
	ret, err := FindMessage(exec, o.ID, o.Date)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *MessageSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *MessageSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MessageSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty MessageSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MessageSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	messages := MessageSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), messagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT `message`.* FROM `message` WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, messagePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(messagePrimaryKeyColumns), 1, len(messagePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&messages)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in MessageSlice")
	}

	*o = messages

	return nil
}

// MessageExists checks if the Message row exists.
func MessageExists(exec boil.Executor, id int64, date time.Time) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from `message` where `id`=? AND `date`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id, date)
	}

	row := exec.QueryRow(sql, id, date)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if message exists")
	}

	return exists, nil
}

// MessageExistsG checks if the Message row exists.
func MessageExistsG(id int64, date time.Time) (bool, error) {
	return MessageExists(boil.GetDB(), id, date)
}

// MessageExistsGP checks if the Message row exists. Panics on error.
func MessageExistsGP(id int64, date time.Time) bool {
	e, err := MessageExists(boil.GetDB(), id, date)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// MessageExistsP checks if the Message row exists. Panics on error.
func MessageExistsP(exec boil.Executor, id int64, date time.Time) bool {
	e, err := MessageExists(exec, id, date)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}


