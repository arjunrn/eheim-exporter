package processor

import (
	"sync"
)

type FilterIDTracker interface {
	Add(string)
	List() []string
}

type filterIDTracker struct {
	lock sync.RWMutex
	ids  map[string]struct{}
}

func (f *filterIDTracker) Add(id string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.ids[id] = struct{}{}
}

func (f *filterIDTracker) List() []string {
	f.lock.RLock()
	defer f.lock.RUnlock()
	results := make([]string, len(f.ids))
	var count int
	for id := range f.ids {
		results[count] = id
		count++
	}
	return results
}

func NewFilterIDTracker() FilterIDTracker {
	return &filterIDTracker{
		ids: make(map[string]struct{}),
	}
}
