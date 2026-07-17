class SlidingWindow:
    """Allows at most `limit` requests in any window_ms span, measured
    from each request.

    TODO: this is the fixed-window counter this exercise exists to
    replace -- it resets at boundaries, so a burst on each side of one
    puts 2x the limit through.
    """

    def __init__(self, limit: int, window_ms: int):
        self.limit = limit
        self.window_ms = window_ms
        self.window_start = 0
        self.count = 0

    def allow_at(self, now_ms: int) -> bool:
        if now_ms - self.window_start >= self.window_ms:
            self.window_start = now_ms
            self.count = 0
        if self.count < self.limit:
            self.count += 1
            return True
        return False
