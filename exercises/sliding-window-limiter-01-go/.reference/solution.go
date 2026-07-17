package main

// SlidingWindow keeps the timestamps of allowed requests and evicts
// the ones a full window old before deciding -- the window slides with
// every request instead of snapping to boundaries. Denied requests are
// never recorded, so they can't consume the budget.
type SlidingWindow struct {
	limit        int
	windowMillis int64
	allowed      []int64
}

func NewSlidingWindow(limit int, windowMillis int64) *SlidingWindow {
	return &SlidingWindow{limit: limit, windowMillis: windowMillis}
}

func (s *SlidingWindow) AllowAt(nowMillis int64) bool {
	cutoff := nowMillis - s.windowMillis
	keep := s.allowed[:0]
	for _, t := range s.allowed {
		if t > cutoff {
			keep = append(keep, t)
		}
	}
	s.allowed = keep

	if len(s.allowed) < s.limit {
		s.allowed = append(s.allowed, nowMillis)
		return true
	}
	return false
}
