package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testAliases(t *testing.T) {
	t.Parallel()

	query := Aliases(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testAliasesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = alia.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAliasesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Aliases(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAliasesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := AliaSlice{alia}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAliasesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := AliaExists(tx, alia.ID)
	if err != nil {
		t.Errorf("Unable to check if Alia exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AliaExistsG to return true, but got false.")
	}
}

func testAliasesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	aliaFound, err := FindAlia(tx, alia.ID)
	if err != nil {
		t.Error(err)
	}

	if aliaFound == nil {
		t.Error("want a record, got nil")
	}
}

func testAliasesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Aliases(tx).Bind(alia); err != nil {
		t.Error(err)
	}
}

func testAliasesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Aliases(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAliasesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	aliaOne := &Alia{}
	aliaTwo := &Alia{}
	if err = randomize.Struct(seed, aliaOne, aliaDBTypes, false, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}
	if err = randomize.Struct(seed, aliaTwo, aliaDBTypes, false, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = aliaOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = aliaTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Aliases(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAliasesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	aliaOne := &Alia{}
	aliaTwo := &Alia{}
	if err = randomize.Struct(seed, aliaOne, aliaDBTypes, false, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}
	if err = randomize.Struct(seed, aliaTwo, aliaDBTypes, false, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = aliaOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = aliaTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func aliaBeforeInsertHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func aliaAfterInsertHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func aliaAfterSelectHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func aliaBeforeUpdateHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func aliaAfterUpdateHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func aliaBeforeDeleteHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func aliaAfterDeleteHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func aliaBeforeUpsertHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func aliaAfterUpsertHook(e boil.Executor, o *Alia) error {
	*o = Alia{}
	return nil
}

func testAliasesHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Alia{}
	o := &Alia{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, aliaDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Alia object: %s", err)
	}

	AddAliaHook(boil.BeforeInsertHook, aliaBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	aliaBeforeInsertHooks = []AliaHook{}

	AddAliaHook(boil.AfterInsertHook, aliaAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	aliaAfterInsertHooks = []AliaHook{}

	AddAliaHook(boil.AfterSelectHook, aliaAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	aliaAfterSelectHooks = []AliaHook{}

	AddAliaHook(boil.BeforeUpdateHook, aliaBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	aliaBeforeUpdateHooks = []AliaHook{}

	AddAliaHook(boil.AfterUpdateHook, aliaAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	aliaAfterUpdateHooks = []AliaHook{}

	AddAliaHook(boil.BeforeDeleteHook, aliaBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	aliaBeforeDeleteHooks = []AliaHook{}

	AddAliaHook(boil.AfterDeleteHook, aliaAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	aliaAfterDeleteHooks = []AliaHook{}

	AddAliaHook(boil.BeforeUpsertHook, aliaBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	aliaBeforeUpsertHooks = []AliaHook{}

	AddAliaHook(boil.AfterUpsertHook, aliaAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	aliaAfterUpsertHooks = []AliaHook{}
}

func testAliasesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAliasesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx, aliaColumns...); err != nil {
		t.Error(err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}











func testAliasesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = alia.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testAliasesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := AliaSlice{alia}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testAliasesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Aliases(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	aliaDBTypes = map[string]string{"ID": "bigint", "Key": "varchar", "Value": "varchar"}
	_           = bytes.MinRead
)

func testAliasesUpdate(t *testing.T) {
	t.Parallel()

	if len(aliaColumns) == len(aliaPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	if err = alia.Update(tx); err != nil {
		t.Error(err)
	}
}

func testAliasesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(aliaColumns) == len(aliaPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	alia := &Alia{}
	if err = randomize.Struct(seed, alia, aliaDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, alia, aliaDBTypes, true, aliaPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(aliaColumns, aliaPrimaryKeyColumns) {
		fields = aliaColumns
	} else {
		fields = strmangle.SetComplement(
			aliaColumns,
			aliaPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(alia))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := AliaSlice{alia}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}

func testAliasesUpsert(t *testing.T) {
	t.Parallel()

	if len(aliaColumns) == len(aliaPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	alia := Alia{}
	if err = randomize.Struct(seed, &alia, aliaDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = alia.Upsert(tx, nil); err != nil {
		t.Errorf("Unable to upsert Alia: %s", err)
	}

	count, err := Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &alia, aliaDBTypes, false, aliaPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Alia struct: %s", err)
	}

	if err = alia.Upsert(tx, nil); err != nil {
		t.Errorf("Unable to upsert Alia: %s", err)
	}

	count, err = Aliases(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

