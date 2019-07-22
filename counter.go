// Package counter provides a library for tracking action times and returning
// their averages.
package counter

import (
	"encoding/json"
	"fmt"
)

// DataStore is an interface for handling new actions.
type DataStore interface {
	// Locks are intended to behave as with sync.RWMutex
	Lock()
	Unlock()
	RLock()
	RUnlock()

	// Add retrieves the given action and adds the given
	// value to its average.
	// See the Average struct for details.
	Add(action string, value int) error

	// Get retrieves all actions and their averages.
	Get() (map[string]*Average, error)
}

// ActionCounter manages adding and retrieving action averages.
type ActionCounter struct {
	DataStore DataStore
}

// AddAction takes marshaled json and adds it to the data store.
func (ac *ActionCounter) AddAction(jsonString string) error {
	var aa actionAddition
	err := json.Unmarshal([]byte(jsonString), &aa)
	if err != nil {
		return fmt.Errorf("error parsing action json: %v", err)
	}
	return ac.add(aa)
}

// GetStats returns the average of each action stat given.
func (ac *ActionCounter) GetStats() string {
	// Locking cannot happen closer to retrieving the map
	// It must stay locked until it finishes unmarshaling,
	// as a simultaneous write could cause a panic.
	//
	// Note that unlock will still happen if the below panic
	// is triggered.
	ac.DataStore.RLock()
	defer ac.DataStore.RUnlock()

	data, err := ac.DataStore.Get()
	if err != nil {
		// Leave error handling to caller.
		// Either the datastore wasn't initialized, or caller is
		// using a custom implementation. If panic is not acceptable,
		// Wrap this call in a function that returns (string, error).
		panic(err)
	}

	// We have no way of handling this error, and unlike the previous case,
	// it isn't worth panicking over. We need to instead have full test coverage.
	resp, _ := json.Marshal(stats(data))
	return string(resp)
}

// add handles the actual adding. It manages data validation and locking.
func (ac *ActionCounter) add(aa actionAddition) error {
	if aa.Time < 1 {
		return fmt.Errorf("non-positive time given to ActionCounter: %d", aa.Time)
	}
	// Unlike reading, writing should lock as close to the write call as possible because
	// no interactions with given data happen after this point.
	ac.DataStore.Lock()
	err := ac.DataStore.Add(aa.Action, aa.Time)
	// Using defer adds overhead unnecessarily when there's only one place to unlock.
	ac.DataStore.Unlock()
	return err
}
