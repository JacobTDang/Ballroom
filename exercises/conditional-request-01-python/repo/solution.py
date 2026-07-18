class ConditionalStore:
    """A versioned key-value store: Get can be conditioned on
    if_none_match, Put on if_match.

    TODO: no versioning at all -- Put always succeeds and overwrites
    unconditionally (the classic lost-update bug this store exists to
    prevent), and Get hands out the same constant etag forever. Every
    rule in the problem statement is still yours to build.
    """

    def __init__(self):
        self.items = {}  # key -> body

    def get(self, key, if_none_match=None):
        if key not in self.items:
            return (404, None, None)
        return (200, "1", self.items[key])

    def put(self, key, body, if_match=None):
        self.items[key] = body
        return (200, "1")

    def delete(self, key):
        self.items.pop(key, None)
