package txslice

import "sync"

type snapshotData[T any] struct {
	mu sync.Mutex

	latest    []*T            // последний snapshot
	versioned map[string][]*T // версионные snapshot'ы

	isAutoLatestSnap bool
}

func (s *snapshotData[T]) setSnapshot(version string, slice []*T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if version == "" {
		s.latest = slice
		return
	}

	s.versioned[version] = slice
}

func (s *snapshotData[T]) getSnapshot(version string) []*T {
	s.mu.Lock()
	defer s.mu.Unlock()

	if version == "" {
		return s.latest
	}

	return s.versioned[version]
}
