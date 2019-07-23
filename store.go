package counter

import (
	"errors"
	"strings"
	"sync"
)

const uninitializedError = "data store was not initialized"

// DefaultDataStore returns an instance of a non-distributed DataStore ready to use.
func DefaultDataStore() DataStore {
	// Also handles the enforcement of store implementing the DataStore interface.
	return &store{
		data:    map[string]*Average{},
		RWMutex: sync.RWMutex{},
	}
}

// store is the default datastore for the action counter.
// It notably converts all alpha characters to lower case,
// does not handle its own lock, and does not validate data.
type store struct {
	data map[string]*Average
	sync.RWMutex
}

// Add adds a value to the count for the given action.
// It will ignore case and assumes that the lock has
// already been obtained by the caller.
func (s *store) Add(action string, value float64) error {
	if s == nil || s.data == nil {
		return errors.New(uninitializedError)
	}

	action = strings.ToLower(action)
	if s.data[action] == nil {
		s.data[action] = new(Average)
	}

	return s.data[action].Add(value)
}

// Get retrieves a copy of the store's data.
// The data is not necessarily safe to modify and assumes that the
// caller obtained the lock.
func (s *store) Get() (map[string]*Average, error) {
	if s == nil || s.data == nil {
		return nil, errors.New(uninitializedError)
	}

	return s.data, nil
}
