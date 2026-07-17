package main

// SlidingWindow allows at most `limit` requests in any windowMillis
// span, measured from each request.
//
// TODO: this is the fixed-window counter this exercise exists to
// replace -- it resets at boundaries, so a burst on each side of one
// puts 2x the limit through.
type SlidingWindow struct {
	limit        int
	windowMillis int64
	windowStart  int64
	count        int
}

func NewSlidingWindow(limit int, windowMillis int64) *SlidingWindow {
	return &SlidingWindow{limit: limit, windowMillis: windowMillis}
}

func (s *SlidingWindow) AllowAt(nowMillis int64) bool {
	if nowMillis-s.windowStart >= s.windowMillis {
		s.windowStart = nowMillis
		s.count = 0
	}
	if s.count < s.limit {
		s.count++
		return true
	}
	return false
}
