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

	err := ac.add(actionAddition{
		Action: "fly",
		Time:   30,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, true, ts.addCalled, "add called")
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

type testStore struct {
	wasLocked, wasUnlocked, wasRLocked, wasRUnlocked bool
	addCalled, getCalled                             bool
	returnError                                      bool
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

func (ts *testStore) Add(_ string, _ int) error {
	ts.addCalled = true
	if ts.returnError {
		return errors.New(innerError)
	}
	return nil
}

func (ts *testStore) Get() (map[string]*Average, error) {
	ts.getCalled = true
	if ts.returnError {
		return nil, errors.New(innerError)
	}
	return ts.store.Get()
}
