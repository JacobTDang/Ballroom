class IdempotencyStore:
    """One record per request key, tracking it through its lifecycle.

    TODO: no fingerprint tracking, no in-flight/completed distinction,
    no deadline at all -- everything after the first Begin just
    replays whatever was last stored, even if nothing ever completed.
    Every rule in the problem statement is still yours to build.
    """

    def __init__(self, ttl_ms):
        self.ttl_ms = ttl_ms
        self.seen = {}  # key -> response, or None if not completed yet

    def begin_at(self, key, fingerprint, now_ms):
        if key not in self.seen:
            self.seen[key] = None
            return ("execute", None)
        return ("replay", self.seen[key])

    def complete_at(self, key, response, now_ms):
        self.seen[key] = response
