package coffeezone

import (
	"testing"
)

func TestPoolAvailable(t *testing.T) {
	poolAvail, expected := DatabasePoolAvailable(), false
	if poolAvail != expected {
		t.Fatalf("DatabasePoolAvailable() returned %t, not %t", poolAvail, expected)
	}

	NewDatabasePool()
	expected = true
	poolAvail = DatabasePoolAvailable()
	if poolAvail != expected {
		t.Fatalf("DatabasePoolAvailable() returned %t, not %t", poolAvail, expected)
	}
}
