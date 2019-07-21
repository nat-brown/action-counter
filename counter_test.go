package counter

import (
	"errors"
	"sync"
	"testing"
)

const innerError = "an error was ordered"

func TestCounterAddActionInvalidJSON(t *testing.T) {
	ts := newTestStore()
	ac := ActionCounter{
		DataStore: &ts,
	}

	err := ac.AddAction(`{"action:"jump","time":30}`)
	if err == nil {
		t.Fatalf("got nil instead of expected error")
	}

	assertEqual(t, "error parsing action json: invalid character 'j' after object key",
		err.Error(), "error message")
	// Ensure we didn't waste time locking.
	assertEqual(t, false, ts.wasLocked, "was locked")
}

func TestCounterAddInvalidTime(t *testing.T) {
	ts := newTestStore()
	ac := ActionCounter{
		DataStore: &ts,
	}

	err := ac.add(actionAddition{
		Action: "fly",
		Time:   -30,
	})
	if err == nil {
		t.Fatalf("got nil instead of expected error")
	}

	assertEqual(t, "non-positive time given to ActionCounter: -30", err.Error(), "error message")
	// Ensure we didn't waste time locking.
	assertEqual(t, false, ts.wasLocked, "was locked")
}

func TestCounterAddError(t *testing.T) {
	ts := newTestStore()
	ac := ActionCounter{
		DataStore: &ts,
	}
	ts.returnError = true

	err := ac.add(actionAddition{
		Action: "fly",
		Time:   30,
	})
	if err == nil {
		t.Fatalf("got nil instead of expected error")
	}

	assertEqual(t, innerError, err.Error(), "error message")
	assertEqual(t, true, ts.wasUnlocked, "was unlocked")
}

func TestCounterAdd(t *testing.T) {
	ts := newTestStore()
	ac := ActionCounter{
		DataStore: &ts,
	}

	const (
		action = "fly"
		time   = 30
	)
	err := ac.add(actionAddition{
		Action: action,
		Time:   time,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, ts.store.data[action].Value(), float64(time), "average")
	assertEqual(t, true, ts.wasLocked, "was locked")
	assertEqual(t, true, ts.wasUnlocked, "was unlocked")
}

func TestGetStats(t *testing.T) {
	ts := newTestStore()
	ac := ActionCounter{
		DataStore: &ts,
	}

	ss := ac.GetStats()
	assertEqual(t, `[{"action":"swim","avg":31}]`, ss, "stats")
	assertEqual(t, true, ts.wasRLocked, "was read locked")
	assertEqual(t, true, ts.wasRUnlocked, "was read unlocked")
}

func TestGetStatsPanic(t *testing.T) {
	ts := newTestStore()
	ac := ActionCounter{
		DataStore: &ts,
	}
	ts.returnError = true

	defer func() {
		if r := recover(); r != nil {
			assertEqual(t, true, ts.wasRLocked, "was read locked")
			assertEqual(t, true, ts.wasRUnlocked, "was read unlocked")
		} else {
			t.Fatal("did not panic")
		}
	}()

	ac.GetStats()
}

func newTestStore() testStore {
	return testStore{store: &store{
		data: map[string]*Average{
			"swim": &Average{
				count: 70,
				value: 30.88,
			},
		},
		RWMutex: sync.RWMutex{},
	}}
}

// testStore uses the DataStore interface to allow introspection into
// what ActionCounter is calling.
type testStore struct {
	// Metadata on called functions.
	wasLocked, wasUnlocked, wasRLocked, wasRUnlocked bool

	// Tells testStore to return errors for subsequent calls.
	returnError bool

	// Gives access to map for Get call.
	*store
}

func (ts *testStore) Lock() {
	ts.wasLocked = true
}

func (ts *testStore) Unlock() {
	ts.wasUnlocked = true
}

func (ts *testStore) RLock() {
	ts.wasRLocked = true
}

func (ts *testStore) RUnlock() {
	ts.wasRUnlocked = true
}

// Add allows testing that the data store errors bubble up.
func (ts *testStore) Add(action string, value int) error {
	if ts.returnError {
		return errors.New(innerError)
	}
	return ts.store.Add(action, value)
}

// Get allows testing that the data store errors bubble up.
func (ts *testStore) Get() (map[string]*Average, error) {
	if ts.returnError {
		return nil, errors.New(innerError)
	}
	return ts.store.Get()
}
