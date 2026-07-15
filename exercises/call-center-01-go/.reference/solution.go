package main

// CallCenter has respondents, managers, and directors. Calls escalate
// respondent -> manager -> director, then queue FIFO.
type CallCenter struct {
	free   map[string]int
	active map[int]string
	queue  []int
}

var levels = []string{"respondent", "manager", "director"}

func NewCallCenter(respondents, managers, directors int) *CallCenter {
	return &CallCenter{
		free: map[string]int{
			"respondent": respondents,
			"manager":    managers,
			"director":   directors,
		},
		active: map[int]string{},
	}
}

// Dispatch routes a new call. Returns the handling level ("respondent",
// "manager", "director") or "queued" when everyone is busy.
func (c *CallCenter) Dispatch(callID int) string {
	for _, level := range levels {
		if c.free[level] > 0 {
			c.free[level]--
			c.active[callID] = level
			return level
		}
	}
	c.queue = append(c.queue, callID)
	return "queued"
}

// EndCall finishes an active call (freeing its employee -- who takes the
// longest-waiting queued call) or abandons a queued one. Returns false
// for unknown/already-ended calls.
func (c *CallCenter) EndCall(callID int) bool {
	if level, ok := c.active[callID]; ok {
		delete(c.active, callID)
		if len(c.queue) > 0 {
			next := c.queue[0]
			c.queue = c.queue[1:]
			c.active[next] = level
		} else {
			c.free[level]++
		}
		return true
	}
	for i, id := range c.queue {
		if id == callID {
			c.queue = append(c.queue[:i], c.queue[i+1:]...)
			return true
		}
	}
	return false
}

// HandlerOf returns the level handling the call, "queued" if waiting,
// or "" if unknown or ended.
func (c *CallCenter) HandlerOf(callID int) string {
	if level, ok := c.active[callID]; ok {
		return level
	}
	for _, id := range c.queue {
		if id == callID {
			return "queued"
		}
	}
	return ""
}
