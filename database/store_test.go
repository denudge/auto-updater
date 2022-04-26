package database

import (
	"github.com/denudge/auto-updater/catalog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatDatabaseStoreIsAStoreInterface(t *testing.T) {
	var store catalog.StoreInterface = nil

	store = &DbCatalogStore{}

	assert.NotNil(t, store)
}
