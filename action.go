package counter

import "encoding/json"

type actionAddition struct {
	Action string `json:"action"`
	Time   int    `json:"time"`
}

type stats map[string]*Average

func (ss stats) MarshalJSON() ([]byte, error) {
	list := make([]stat, len(ss))
	var i int
	for action, avg := range ss {
		list[i] = stat{
			Action:  action,
			Average: avg.IntValue(),
		}
		i++
	}
	return json.Marshal(list)
}

type stat struct {
	Action  string `json:"action"`
	Average int    `json:"avg"`
}
