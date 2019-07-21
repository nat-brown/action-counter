package counter

import (
	"fmt"
	"math"
	"math/bits"
	"testing"
)

const requiredPercentAccuracy = 0.99999 // 99.999% ("five nines")

// assertAccurate checks that the actual value differs from the expected
// less than the required percentage.
func assertAccurate(t *testing.T, expected, actual float64, failMsg string) {
	if (actual / expected) < requiredPercentAccuracy {
		t.Fatalf(failMsg)
	}
}

func TestAverageAddMany(t *testing.T) {
	cases := []struct {
		value           int
		expectedAverage float64
	}{
		{value: 10, expectedAverage: 10},
		{value: 5, expectedAverage: 7.5},
		{value: 8, expectedAverage: 23.0 / 3},
	}

	average := new(Average)
	for i, c := range cases {
		err := average.Add(c.value)
		if err != nil {
			t.Fatalf("case %d failed: unexpected error: %v", i, err)
		}

		assertAccurate(t, c.expectedAverage, average.value,
			fmt.Sprintf(`case %d failed:
	adding value: %d
	expected average: %f 
	actual average: %f`,
				i, c.value, c.expectedAverage, average.value),
		)
	}
}

func TestAverageAccuracy(t *testing.T) {
	// Averaging the same value repeatedly will surface any noticeable
	// loss in accuracy from imprecise float representation.
	avg := new(Average)
	const (
		val            = 5555
		floatVal       = float64(val)
		requiredRounds = 1e6 // 1 million actions
	)

	for i := 0; i < requiredRounds; i++ {
		err := avg.Add(val)
		if err != nil {
			t.Fatalf("case %d failed: unexpected error %v", i, err)
		}
		assertAccurate(t, floatVal, avg.value,
			fmt.Sprintf("Loss of precision exceeded acceptable values by %d round: %f",
				i, (avg.value/floatVal)),
		)
	}
}

// maxInt calculates the maximum integer value regardless of
// the operating system.
func maxInt() int {
	// Operations on signed integers behave differently from unsigned integers,
	// so we will use unsigned and convert at the end.
	// The end of this section discusses the two operators used
	// in this function: https://golang.org/ref/spec#Arithmetic_operators
	//
	// Relies on the fact that uint and int are the same number of bits,
	// differing by a sign flag on the leftmost bit for int.
	// See: https://golang.org/ref/spec#Numeric_types
	allOnes := ^uint(0)

	// Leftmost bit of 1 will indicate a negative, so we need to change
	// it to 0. Shift operation fills "empty" space with 0.
	toggledSign := allOnes >> 1

	return int(toggledSign)
}

func TestMaxInt(t *testing.T) {
	actual := maxInt()
	var expected int64 = math.MaxInt32
	if bits.UintSize == 64 {
		// This line will not compile on a 32 bit system if the type of
		// "expected" is not int64. It will compile without error on a
		// 64 bit system if type is int.
		expected = math.MaxInt64
	}
	if int64(actual) != expected {
		t.Fatalf("Maxint (%d) was not correctly calculated for Uintsize %d", actual, bits.UintSize)
	}
}

func TestAverageEdgeCases(t *testing.T) {
	// The expected values here were calculated using Wolfram Alpha,
	// as most calculators will display less precision.
	tts := []struct {
		name     string
		avg      Average
		addValue int

		shouldErr   bool
		errorMsg    string
		expectedAvg Average
	}{
		{
			name: "start max average",
			avg: Average{
				count: 11,
				value: float64(math.MaxInt64),
			},
			addValue: 1,
			expectedAvg: Average{
				count: 12,
				value: 8454757700450211157 + 5.0/12,
			},
		}, {
			name: "start max count",
			avg: Average{
				count: maxInt(),
				value: 10,
			},
			addValue:  1,
			shouldErr: true,
			errorMsg:  "Average.count overflow - cannot add to this Average",
		}, {
			name: "zero addition",
			avg: Average{
				count: 5,
				value: 10,
			},
			addValue:  0,
			shouldErr: true,
			errorMsg:  "Average.Add called with non-positive value 0",
		}, {
			name: "negative addition",
			avg: Average{
				count: 5,
				value: 10,
			},
			addValue:  -1,
			shouldErr: true,
			errorMsg:  "Average.Add called with non-positive value -1",
		}, {
			name: "max addition",
			avg: Average{
				count: 500,
				value: 10,
			},
			addValue: maxInt(),
			expectedAvg: Average{
				count: 501,
				value: 18409924225259043 + 265.0/501,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.avg.Add(tt.addValue)
			if !tt.shouldErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.shouldErr {
				assertEqual(t, tt.errorMsg, err.Error(), "error message")
				return
			}
			assertAccurate(t, tt.expectedAvg.value, tt.avg.value,
				fmt.Sprintf("\nexpected average: %+v\nactual average: %+v",
					tt.expectedAvg, tt.avg))
		})
	}
}

func TestAverageHelpers(t *testing.T) {
	const (
		val = 8
		ct  = 30
	)
	actual := NewAverage(val, ct)
	expected := &Average{
		value: val,
		count: ct,
	}
	assertEqual(t, *expected, *actual, "average")
	assertEqual(t, actual.Value(), actual.value, "value")
	assertEqual(t, actual.Count(), actual.count, "count")
}

func TestIntAverage(t *testing.T) {
	tts := []struct {
		value    float64
		expected int
	}{
		{value: 0, expected: 0},
		{value: 10.1, expected: 10},
		{value: 10.5, expected: 11},
	}

	for _, tt := range tts {
		t.Run(fmt.Sprintf("%f rounds to %d", tt.value, tt.expected), func(t *testing.T) {
			avg := Average{value: tt.value}
			assertEqual(t, tt.expected, avg.IntValue(), "rounded value")
		})
	}
}
