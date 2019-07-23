# Action-Counter

This library provides a simple tool for tracking actions and their average times in a concurrency-safe manner. It provides two functions:
* Add an action
* Retrieve current averages

Using it would look similar to below:
```
package anything

import (
    "fmt"
    counter "github.com/nat-brown/action-counter"
)

func CountSomeActions() {
    ac := counter.ActionCounter{
        DataStore: counter.DefaultDataStore(),
    }

    ac.AddAction(`{"action":"jump", "time":105}`)
    ac.AddAction(`{"action":"run", "time":75.3}`)
    ac.AddAction(`{"action":"jump", "time":200}`)

    output = ac.GetStats()
    fmt.Println(output) // [{"action":"jump","avg":152.5},{"action":"run","avg":75.3}]
}
```
---
## Interface

The action counter works exclusively in serialized json strings.

#### Adding an action

Action additions expect a json object containing two attributes: `action` and `time`. `time` must be a positive and non-zero number, while `action` is a case-insensitive string and will always be stored as the lowercase version. Extra key-value pairs will be ignored.

If your project requires case-sensitivity or other action name/time specifications, you can define your own object implementing the `DataStore` interface.

Note that `time > 0` is a requirement handled outside of `DataStore`. Any limitations added on the `DataStore` level will stack with that requirement.

#### Retrieving statistics

Retrieved statistics will be a list of json objects containing an `action` and its corresponding average, `avg`. `action` will be a string and `avg` will be a json `number`, which may be a decimal and should be unmarshaled as a float64.

#### The DataStore interface

While this library works concurrently, it may not work for your needs with the default DataStore, such as in the case of a distributed system. If you need custom data storage and retrieval, or different constraints on input, this interface makes room for that.

The `Average` struct has public functions to facilitate retrieving data from other sources, instantiating an `Average` instance, letting it take care of the averaging logic, and then pushing the result back to a database.

---
## Installation

This package assumes that you have a working Go environment. It has been tested in go version 1.11, but does not use dependencies that would likely prohibit earlier versions. If you use an earlier version of go, be sure to run the tests to ensure your version is compatible with this library.

To install the `counter` library, run the following `go get` command

```
go get github.com/nat-brown/action-counter
```

It can now be imported into go code.

```
import (
    counter "github.com/nat-brown/action-counter"
)
```

Note that the project code will be available at `$GOPATH/src/github.com/nat-brown/action-counter` on your system. For more details (such as the case for multiple GOPATHs), you can read details [here](https://golang.org/pkg/cmd/go/internal/get/).

---
## Running Tests

Tests can be run from the `action-counter` directory with the `go test` command, or `go test -v` for a list of all tests being run. For more options on running tests, you can read more [here](https://golang.org/pkg/cmd/go/internal/test/).

---
## Go Docs

This repository was written with godoc in mind. To run godoc, you can read instructions [here](https://godoc.org/golang.org/x/tools/cmd/godoc).