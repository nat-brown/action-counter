package counter

import (
	"encoding/json"
	"testing"
)

func TestStatsMarshal(t *testing.T) {
	tts := []struct {
		name     string
		data     map[string]*Average
		expected string
	}{
		{
			name: "filled data",
			data: map[string]*Average{
				"jump": &Average{value: 5.2},
				"run":  &Average{value: 3.6},
			},
			expected: `[{"action":"jump","avg":5},{"action":"run","avg":4}]`,
		}, {
			name:     "empty map",
			data:     map[string]*Average{},
			expected: "[]",
		}, {
			name:     "nil map",
			expected: "[]",
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := json.Marshal(stats(tt.data))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Convert to string because slices aren't inherently comparable.
			assertEqual(t, tt.expected, string(actual), "marshaled data")
		})
	}
}
