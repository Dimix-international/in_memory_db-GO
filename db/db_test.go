package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShardDB(t *testing.T) {
	db := NewShardMap()

	//Set
	db.Set("a", "1")
	db.Set("b", "2")
	db.Set("c", "3")

	//Get
	value, ok := db.Get("a")

	assert.True(t, ok)
	assert.Equal(t, value, "1")

	//Delete
	db.Delete("b")
	value, ok = db.Get("b")

	assert.Equal(t, ok, false)
	assert.Equal(t, value, "")

	//Set like update
	db.Set("c", "4")
	value, ok = db.Get("c")

	assert.True(t, ok)
	assert.Equal(t, value, "4")
}
