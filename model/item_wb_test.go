package model

import (
	"testing"
)

func TestNewItem(t *testing.T) {
	// when
	item := NewItem()
	// then
	if item.State != Newly {
		t.Error("expect state to be Newly")
	} else if item.Type != ITEM {
		t.Error("expect type to be ITEM")
	}
}
