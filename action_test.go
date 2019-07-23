package counter

import (
	"encoding/json"
	"testing"
)

func TestStatsMarshal(t *testing.T) {
	tts := []struct {
		name     string
		data     map[string]*Average
		expected []string // Multiple potential responses because maps aren't ordered.
	}{
		{
			name: "filled data",
			data: map[string]*Average{
				"jump": &Average{value: 5.2},
				"run":  &Average{value: 3.6},
			},
			expected: []string{
				`[{"action":"jump","avg":5.2},{"action":"run","avg":3.6}]`,
				`[{"action":"run","avg":3.6},{"action":"jump","avg":5.2}]`,
			},
		}, {
			name:     "empty map",
			data:     map[string]*Average{},
			expected: []string{"[]"},
		}, {
			name:     "nil map",
			expected: []string{"[]"},
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := json.Marshal(stats(tt.data))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var matched bool
			for _, match := range tt.expected {
				// Convert to string because slices aren't inherently comparable.
				if string(actual) == match {
					matched = true
				}
			}
			if !matched {
				t.Fatalf("\nexpected any of: %v\nactual: %s", tt.expected, actual)
			}
		})
	}
}
