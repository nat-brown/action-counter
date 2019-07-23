package counter

import (
	"os"
	"testing"
)

// assertEqual is a helper for the most common test assertion.
func assertEqual(t *testing.T, expected, actual interface{}, name string) {
	if expected != actual {
		t.Fatalf("\nexpected %s: %+v\nactual %s: %+v", name, expected, name, actual)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
