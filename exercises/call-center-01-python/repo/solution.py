class CallCenter:
    """Call center with respondents, managers, and directors. Calls
    escalate respondent -> manager -> director, then queue FIFO."""

    def __init__(self, respondents: int, managers: int, directors: int):
        pass

    def dispatch(self, call_id: int) -> str:
        """Route a new call. Return the handling level ("respondent",
        "manager", "director") or "queued" when everyone is busy."""
        raise NotImplementedError

    def end_call(self, call_id: int) -> bool:
        """Finish an active call (freeing its employee -- who takes the
        longest-waiting queued call) or abandon a queued one. Return
        False for unknown/already-ended calls."""
        raise NotImplementedError

    def handler_of(self, call_id: int) -> str:
        """The level handling the call, "queued" if waiting, or "" if
        unknown or ended."""
        raise NotImplementedError
