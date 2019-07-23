package counter_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	counter "github.com/nat-brown/action-counter"
)

func Example() {
	ac := counter.ActionCounter{
		DataStore: counter.DefaultDataStore(),
	}

	type action struct {
		action string
		time   int
	}
	actions := []action{
		{action: "jump", time: 30},
		{action: "run", time: 3600},
		{action: "run", time: 3300},
		{action: "jump", time: 30},
		{action: "run", time: 3337},
		{action: "swim", time: 250},
		{action: "jump", time: 32},
	}

	var wg sync.WaitGroup
	wg.Add(len(actions))
	for _, a := range actions {
		go func(a action) {
			request := fmt.Sprintf(`{"action":"%s","time":%d}`, a.action, a.time)
			err := ac.AddAction(request)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(a)
	}
	wg.Wait()

	output := ac.GetStats()
	fmt.Println(sortKeys(output))

	// Output:
	// [{"action":"jump","avg":30.666666666666664},{"action":"run","avg":3412.333333333333},{"action":"swim","avg":250}]
}

func ExampleActionCounter() {
	ac := counter.ActionCounter{
		DataStore: counter.DefaultDataStore(),
	}

	output := ac.GetStats()
	fmt.Println(output)

	ac.AddAction(`{"action":"jump", "time":105}`)
	ac.AddAction(`{"action":"run", "time":75.3}`)
	ac.AddAction(`{"action":"jump", "time":200}`)

	output = ac.GetStats()
	fmt.Println(sortKeys(output))

	// Output:
	// []
	// [{"action":"jump","avg":152.5},{"action":"run","avg":75.3}]
}

func sortKeys(output string) string {
	averages := []struct {
		Action string  `json:"action"`
		Avg    float64 `json:"avg"`
	}{}
	json.Unmarshal([]byte(output), &averages)

	sort.Slice(averages, func(i, j int) bool {
		return averages[i].Action < averages[j].Action
	})

	b, _ := json.Marshal(averages)
	return string(b)
}
