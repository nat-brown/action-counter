package counter

import (
	"errors"
	"fmt"
)

// NewAverage initializes a new Average by setting private values.
// This function is intended for instantiating an Average using data
// from an outside data store.
//
// Note that new(Average) is equivalent to NewAverage(0, 0).
func NewAverage(value float64, count int) *Average {
	return &Average{
		value: value,
		count: count,
	}
}

// Average represents an action's cumulative average.
type Average struct {
	count int
	value float64
}

// Add adds a new value to the cumulative average.
func (a *Average) Add(newVal float64) error {
	if newVal <= 0 {
		return fmt.Errorf("Average.Add called with non-positive value %f", newVal)
	}
	if a.count+1 < 0 {
		return errors.New("Average.count overflow - cannot add to this Average")
	}

	newCount := float64(a.count + 1)
	largerRatio := (newCount - 1) / newCount
	smallerRatio := float64(1) / newCount

	newAverage := largerRatio*a.value + smallerRatio*newVal

	a.value = newAverage
	a.count++

	return nil
}

// Count returns the number of data points in the calculated average.
func (a *Average) Count() int { return a.count }

// Value returns the calculated average value.
func (a *Average) Value() float64 { return a.value }
