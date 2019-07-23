package counter

import "encoding/json"

// actionAddition is the action "request" struct.
type actionAddition struct {
	Action string  `json:"action"`
	Time   float64 `json:"time"`
}

// stats aliases a map of an action to an average to enable custom marshaling.
type stats map[string]*Average

// MarhsalJSON converts stats to a list of stat structs.
func (ss stats) MarshalJSON() ([]byte, error) {
	list := make([]stat, len(ss))
	var i int
	for action, avg := range ss {
		list[i] = stat{
			Action:  action,
			Average: avg.Value(),
		}
		i++
	}
	return json.Marshal(list)
}

// stat is the action average "response" struct.
type stat struct {
	Action  string  `json:"action"`
	Average float64 `json:"avg"`
}
