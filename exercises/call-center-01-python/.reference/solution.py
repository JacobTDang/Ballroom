_LEVELS = ["respondent", "manager", "director"]


class CallCenter:
    """Call center with respondents, managers, and directors. Calls
    escalate respondent -> manager -> director, then queue FIFO."""

    def __init__(self, respondents: int, managers: int, directors: int):
        self._free = {
            "respondent": respondents,
            "manager": managers,
            "director": directors,
        }
        self._active = {}  # call_id -> level
        self._queue = []   # call_ids in arrival order

    def dispatch(self, call_id: int) -> str:
        for level in _LEVELS:
            if self._free[level] > 0:
                self._free[level] -= 1
                self._active[call_id] = level
                return level
        self._queue.append(call_id)
        return "queued"

    def end_call(self, call_id: int) -> bool:
        if call_id in self._active:
            level = self._active.pop(call_id)
            if self._queue:
                nxt = self._queue.pop(0)
                self._active[nxt] = level
            else:
                self._free[level] += 1
            return True
        if call_id in self._queue:
            self._queue.remove(call_id)
            return True
        return False

    def handler_of(self, call_id: int) -> str:
        if call_id in self._active:
            return self._active[call_id]
        if call_id in self._queue:
            return "queued"
        return ""
