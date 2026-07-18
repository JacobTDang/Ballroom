class IdempotencyStore:
    """One record per key: fingerprint, status (in-flight or
    completed), a deadline that Begin sets and Complete renews, and
    (once completed) the stored response. Past its deadline, a record
    is treated as if it never existed -- Begin starts clean, and
    Complete has nothing to attach to."""

    EXECUTE = "execute"
    IN_FLIGHT = "in-flight"
    REPLAY = "replay"

    def __init__(self, ttl_ms):
        self.ttl_ms = ttl_ms
        self.records = {}  # key -> {fingerprint, status, deadline, response}

    def begin_at(self, key, fingerprint, now_ms):
        record = self.records.get(key)
        if record is None or now_ms >= record["deadline"]:
            self.records[key] = {
                "fingerprint": fingerprint,
                "status": "in-flight",
                "deadline": now_ms + self.ttl_ms,
                "response": None,
            }
            return (self.EXECUTE, None)

        if record["fingerprint"] != fingerprint:
            raise ValueError(f"fingerprint conflict for key {key!r}")

        if record["status"] == "in-flight":
            return (self.IN_FLIGHT, None)
        return (self.REPLAY, record["response"])

    def complete_at(self, key, response, now_ms):
        record = self.records.get(key)
        if record is None or now_ms >= record["deadline"] or record["status"] != "in-flight":
            raise ValueError(f"no in-flight request for key {key!r}")
        record["status"] = "completed"
        record["response"] = response
        record["deadline"] = now_ms + self.ttl_ms
