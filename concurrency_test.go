package counter

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Since maps are not safe to lock on a per-key basis,
// it's sufficient to test this using only one action type.
// Other tests handle asserting that Action Counter can
// handle multiple actions.
const action = "action"

const (
	add = iota
	read
)

// Choice is for tracking what was added and retrieved
// from the Action Counter.
type choice struct {
	action int
	number float64
}

func TestCounterConcurrency(t *testing.T) {
	dataChannel := make(chan choice)
	allSent := make(chan bool)

	go callCounter(t, dataChannel, allSent)
	results := gatherResults(dataChannel, allSent)
	checkResults(t, results)
}

// callCounter calls an action counter, outputting data to dataChannel and signaling
// when all data has been sent.
func callCounter(t *testing.T, dataChannel chan<- choice, allSent chan<- bool) {
	ac := ActionCounter{
		DataStore: &recorderStore{
			ch: dataChannel,
			store: &store{
				data:    map[string]*Average{},
				RWMutex: sync.RWMutex{},
			},
		},
	}
	var wg sync.WaitGroup

	makeRandomCounterCalls(t, &wg, ac)

	wg.Wait()
	allSent <- true
	close(dataChannel)
	close(allSent)
}

// makeRandomCounterCalls will make a random call for the given action counter until it times out.
func makeRandomCounterCalls(t *testing.T, wg *sync.WaitGroup, ac ActionCounter) {
	timeout := time.After(1 * time.Second)
	for {
		select {
		case <-timeout:
			return
		default:
			go func() {
				wg.Add(1)
				switch rand.Int() % 2 {
				case 0:
					ac.GetStats()
				default: // default instead of 1 to make changing the ratio easier.
					err := ac.AddAction(fmt.Sprintf(`{"action":"%s","time":%d}`, action, rand.Int()))
					if err != nil {
						t.Fatal(err)
					}
				}
				wg.Done()
			}()
		}
	}
}

// GatherResults handles the streamed results from the calls to the action counter.
func gatherResults(dataChannel <-chan choice, allSent <-chan bool) []choice {
	results := []choice{}
	for {
		select {
		case v := <-dataChannel:
			results = append(results, v)
		case <-allSent:
			return results
		}
	}
}

// CheckResults checks that the collected data matches as if it had been called
// sequentially. It relies on other tests for sequential correctness.
func checkResults(t *testing.T, results []choice) {
	ac := ActionCounter{
		DataStore: DefaultDataStore(),
	}
	for _, result := range results {
		switch result.action {
		case read:
			checkRead(t, result, ac)
		case add:
			ac.AddAction(fmt.Sprintf(`{"action":"%s","time":%f}`, action, result.number))
		}
	}
}

// CheckRead handles the required marshaling for asserting that a GetStats call
// matches the expected value.
func checkRead(t *testing.T, result choice, ac ActionCounter) {
	stats := ac.GetStats()
	if result.number == 0 {
		assertEqual(t, "[]", stats, "get stats result")
		return
	}

	b, err := json.Marshal([]struct {
		Action  string  `json:"action"`
		Average float64 `json:"avg"`
	}{{
		Action:  action,
		Average: result.number,
	}})
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, string(b), stats, "get stats result")
}

// RecorderStore tracks data sent over the test store
// by taking advantage of locks.
// We cannot test concurrent adding and reading without locking
// during value recordings (to preserve action order for determining
// expectations), and locking from the outside would mask errors in
// locks from the inside. We thus need to only write expectations
// from within the DataStore's own functions, which would utilize its locks.
type recorderStore struct {
	ch chan<- choice
	*store
}

// Add takes advantage of the fact that Lock() is called before this
// (as tested elsewhere). If Lock() isn't called, data should be mingled
// in incorrect order in the channel.
func (rs *recorderStore) Add(action string, value float64) error {
	rs.pushAdd(value)
	err := rs.store.Add(action, value)

	// Simulate a read, whether or not there was one, to allow us to
	// check the value as we go.
	// This does not handle asserting concurrent reads for an Action
	// Counter, as we've only asserted that AddAction will lock.
	rs.pushRead()
	return err
}

// RUnlock assumes that RLock() was called earlier and does not
// care about preventing mingled read calls from entering the channel
// in random order, since the values pushed to the channel will be identical.
func (rs *recorderStore) RUnlock() {
	rs.pushRead()
	rs.store.RUnlock()
}

func (rs *recorderStore) pushRead() {
	av := rs.store.data[action]
	var val float64
	if av != nil {
		val = av.Value()
	}
	rs.ch <- choice{
		action: read,
		number: val,
	}
}

func (rs *recorderStore) pushAdd(value float64) {
	rs.ch <- choice{
		action: add,
		number: value,
	}
}
