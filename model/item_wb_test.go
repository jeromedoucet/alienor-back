package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
	// when
	item := NewItem()
	// then
	assert.Equal(t, Newly, item.State)
	assert.Equal(t, ITEM, item.Type)
}
