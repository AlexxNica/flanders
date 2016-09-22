package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Aliases", testAliases)
	t.Run("Messages", testMessages)
}

func TestDelete(t *testing.T) {
	t.Run("Aliases", testAliasesDelete)
	t.Run("Messages", testMessagesDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Aliases", testAliasesQueryDeleteAll)
	t.Run("Messages", testMessagesQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Aliases", testAliasesSliceDeleteAll)
	t.Run("Messages", testMessagesSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Aliases", testAliasesExists)
	t.Run("Messages", testMessagesExists)
}

func TestFind(t *testing.T) {
	t.Run("Aliases", testAliasesFind)
	t.Run("Messages", testMessagesFind)
}

func TestBind(t *testing.T) {
	t.Run("Aliases", testAliasesBind)
	t.Run("Messages", testMessagesBind)
}

func TestOne(t *testing.T) {
	t.Run("Aliases", testAliasesOne)
	t.Run("Messages", testMessagesOne)
}

func TestAll(t *testing.T) {
	t.Run("Aliases", testAliasesAll)
	t.Run("Messages", testMessagesAll)
}

func TestCount(t *testing.T) {
	t.Run("Aliases", testAliasesCount)
	t.Run("Messages", testMessagesCount)
}

func TestHooks(t *testing.T) {
	t.Run("Aliases", testAliasesHooks)
	t.Run("Messages", testMessagesHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Aliases", testAliasesInsert)
	t.Run("Aliases", testAliasesInsertWhitelist)
	t.Run("Messages", testMessagesInsert)
	t.Run("Messages", testMessagesInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {}

func TestReload(t *testing.T) {
	t.Run("Aliases", testAliasesReload)
	t.Run("Messages", testMessagesReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Aliases", testAliasesReloadAll)
	t.Run("Messages", testMessagesReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Aliases", testAliasesSelect)
	t.Run("Messages", testMessagesSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Aliases", testAliasesUpdate)
	t.Run("Messages", testMessagesUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Aliases", testAliasesSliceUpdateAll)
	t.Run("Messages", testMessagesSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("Aliases", testAliasesUpsert)
	t.Run("Messages", testMessagesUpsert)
}

