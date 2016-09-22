package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testMessages(t *testing.T) {
	t.Parallel()

	query := Messages(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testMessagesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = message.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMessagesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Messages(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMessagesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := MessageSlice{message}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMessagesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := MessageExists(tx, message.ID, message.Date)
	if err != nil {
		t.Errorf("Unable to check if Message exists: %s", err)
	}
	if !e {
		t.Errorf("Expected MessageExistsG to return true, but got false.")
	}
}

func testMessagesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	messageFound, err := FindMessage(tx, message.ID, message.Date)
	if err != nil {
		t.Error(err)
	}

	if messageFound == nil {
		t.Error("want a record, got nil")
	}
}

func testMessagesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Messages(tx).Bind(message); err != nil {
		t.Error(err)
	}
}

func testMessagesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Messages(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testMessagesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	messageOne := &Message{}
	messageTwo := &Message{}
	if err = randomize.Struct(seed, messageOne, messageDBTypes, false, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}
	if err = randomize.Struct(seed, messageTwo, messageDBTypes, false, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = messageOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = messageTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Messages(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testMessagesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	messageOne := &Message{}
	messageTwo := &Message{}
	if err = randomize.Struct(seed, messageOne, messageDBTypes, false, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}
	if err = randomize.Struct(seed, messageTwo, messageDBTypes, false, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = messageOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = messageTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func messageBeforeInsertHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func messageAfterInsertHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func messageAfterSelectHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func messageBeforeUpdateHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func messageAfterUpdateHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func messageBeforeDeleteHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func messageAfterDeleteHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func messageBeforeUpsertHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func messageAfterUpsertHook(e boil.Executor, o *Message) error {
	*o = Message{}
	return nil
}

func testMessagesHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Message{}
	o := &Message{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, messageDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Message object: %s", err)
	}

	AddMessageHook(boil.BeforeInsertHook, messageBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	messageBeforeInsertHooks = []MessageHook{}

	AddMessageHook(boil.AfterInsertHook, messageAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	messageAfterInsertHooks = []MessageHook{}

	AddMessageHook(boil.AfterSelectHook, messageAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	messageAfterSelectHooks = []MessageHook{}

	AddMessageHook(boil.BeforeUpdateHook, messageBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	messageBeforeUpdateHooks = []MessageHook{}

	AddMessageHook(boil.AfterUpdateHook, messageAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	messageAfterUpdateHooks = []MessageHook{}

	AddMessageHook(boil.BeforeDeleteHook, messageBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	messageBeforeDeleteHooks = []MessageHook{}

	AddMessageHook(boil.AfterDeleteHook, messageAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	messageAfterDeleteHooks = []MessageHook{}

	AddMessageHook(boil.BeforeUpsertHook, messageBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	messageBeforeUpsertHooks = []MessageHook{}

	AddMessageHook(boil.AfterUpsertHook, messageAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	messageAfterUpsertHooks = []MessageHook{}
}

func testMessagesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMessagesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx, messageColumns...); err != nil {
		t.Error(err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}











func testMessagesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = message.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testMessagesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := MessageSlice{message}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testMessagesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Messages(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	messageDBTypes = map[string]string{"Auth": "varchar", "AuthUser": "varchar", "Callid": "varchar", "CallidAleg": "varchar", "ContactIp": "varchar", "ContactPort": "int", "ContactUser": "varchar", "ContentType": "varchar", "CorrelationID": "varchar", "Cseq": "varchar", "CustomField1": "varchar", "CustomField2": "varchar", "CustomField3": "varchar", "Date": "timestamp", "DestinationIp": "varchar", "DestinationPort": "int", "Diversion": "varchar", "Family": "int", "FromDomain": "varchar", "FromTag": "varchar", "FromUser": "varchar", "ID": "bigint", "MSG": "varchar", "Method": "varchar", "MicroTS": "bigint", "Node": "varchar", "OriginatorIp": "varchar", "OriginatorPort": "int", "PidUser": "varchar", "Proto": "int", "RTPStat": "varchar", "Reason": "varchar", "ReplyReason": "varchar", "Ruri": "varchar", "RuriDomain": "varchar", "RuriUser": "varchar", "SourceIp": "varchar", "SourcePort": "int", "ToDomain": "varchar", "ToTag": "varchar", "ToUser": "varchar", "Type": "int", "UserAgent": "varchar", "Via1": "varchar", "Via1Branch": "varchar"}
	_              = bytes.MinRead
)

func testMessagesUpdate(t *testing.T) {
	t.Parallel()

	if len(messageColumns) == len(messagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, message, messageDBTypes, true, messageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	if err = message.Update(tx); err != nil {
		t.Error(err)
	}
}

func testMessagesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(messageColumns) == len(messagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	message := &Message{}
	if err = randomize.Struct(seed, message, messageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, message, messageDBTypes, true, messagePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(messageColumns, messagePrimaryKeyColumns) {
		fields = messageColumns
	} else {
		fields = strmangle.SetComplement(
			messageColumns,
			messagePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(message))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := MessageSlice{message}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}

func testMessagesUpsert(t *testing.T) {
	t.Parallel()

	if len(messageColumns) == len(messagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	message := Message{}
	if err = randomize.Struct(seed, &message, messageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = message.Upsert(tx, nil); err != nil {
		t.Errorf("Unable to upsert Message: %s", err)
	}

	count, err := Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &message, messageDBTypes, false, messagePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Message struct: %s", err)
	}

	if err = message.Upsert(tx, nil); err != nil {
		t.Errorf("Unable to upsert Message: %s", err)
	}

	count, err = Messages(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

