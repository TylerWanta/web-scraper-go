package urlTracker

import (
	"strings"
	"sync"
)

type URLTracker struct {
	sync.Mutex
	data     map[string]bool
	capacity int
}

func NewURlTracker(capacity int) *URLTracker {
	return &URLTracker{
		data:     make(map[string]bool, capacity),
		capacity: capacity,
	}
}

func (sm *URLTracker) Add(s string) bool {
	sm.Lock()
	defer sm.Unlock()

	if len(sm.data) == sm.capacity {
		return false
	}

	if _, ok := sm.data[s]; ok {
		return false
	}

	sm.data[s] = true
	return true
}

func (sm *URLTracker) Length() int {
	sm.Lock()
	defer sm.Unlock()

	return len(sm.data)
}

func (sm *URLTracker) String() string {
	sm.Lock()
	defer sm.Unlock()

	keys := make([]string, 0, len(sm.data))

	for key := range sm.data {
		keys = append(keys, key)
	}

	return strings.Join(keys, "\n")
}
