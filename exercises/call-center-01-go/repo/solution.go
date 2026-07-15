package main

// CallCenter has respondents, managers, and directors. Calls escalate
// respondent -> manager -> director, then queue FIFO.
type CallCenter struct {
}

func NewCallCenter(respondents, managers, directors int) *CallCenter {
	return &CallCenter{}
}

// Dispatch routes a new call. Returns the handling level ("respondent",
// "manager", "director") or "queued" when everyone is busy.
func (c *CallCenter) Dispatch(callID int) string {
	// TODO: implement
	return ""
}

// EndCall finishes an active call (freeing its employee -- who takes the
// longest-waiting queued call) or abandons a queued one. Returns false
// for unknown/already-ended calls.
func (c *CallCenter) EndCall(callID int) bool {
	// TODO: implement
	return false
}

// HandlerOf returns the level handling the call, "queued" if waiting,
// or "" if unknown or ended.
func (c *CallCenter) HandlerOf(callID int) string {
	// TODO: implement
	return ""
}
