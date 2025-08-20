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

	dataCopy := make([]*T, len(slice))
	copy(dataCopy, slice)

	if version == "" {
		s.latest = dataCopy
		return
	}

	if s.versioned == nil {
		s.versioned = make(map[string][]*T)
	}

	s.versioned[version] = dataCopy
}

func (s *snapshotData[T]) getSnapshot(version string) []*T {
	s.mu.Lock()
	defer s.mu.Unlock()

	var snap []*T

	if version == "" {
		snap = s.latest
	} else {
		snap = s.versioned[version]
	}

	if snap == nil {
		return nil
	}

	copySnap := make([]*T, len(snap))
	copy(copySnap, snap)

	return copySnap
}
