package counter

import (
	"sync"
	"testing"
)

func TestStoreHappyPath(t *testing.T) {
	s := store{
		data:    map[string]*Average{},
		RWMutex: sync.RWMutex{},
	}

	err := s.Add("jump", 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = s.Add("Jump", 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := s.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, float64(25), data["jump"].Value(), "average")
}

func TestNilHandling(t *testing.T) {
	var s store
	err := s.Add("anything", 5)
	assertEqual(t, uninitializedError, err.Error(), "error")
	_, err = s.Get()
	assertEqual(t, uninitializedError, err.Error(), "error")
}
