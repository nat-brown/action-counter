package counter

import "testing"

func assertEqual(t *testing.T, expected, actual interface{}, name string) {
	if expected != actual {
		t.Fatalf("\nexpected %s: %+v\nactual %s: %+v", expected, name, actual, name)
	}
}
