class ConditionalStore:
    """Every successful write (create or update) draws its etag from
    one store-wide monotonic sequence -- so deleting a key and
    recreating it never reissues an old etag, and a stale cached etag
    can never falsely match again, on either the read or write side."""

    def __init__(self):
        self.items = {}  # key -> {"etag": str, "body": ...}
        self._seq = 0

    def _next_etag(self):
        self._seq += 1
        return str(self._seq)

    def get(self, key, if_none_match=None):
        entry = self.items.get(key)
        if entry is None:
            return (404, None, None)
        if if_none_match and if_none_match == entry["etag"]:
            return (304, None, None)
        return (200, entry["etag"], entry["body"])

    def put(self, key, body, if_match=None):
        entry = self.items.get(key)
        if entry is None:
            if if_match:
                # Claimed a version of something that doesn't exist.
                return (412, None)
            etag = self._next_etag()
            self.items[key] = {"etag": etag, "body": body}
            return (200, etag)

        if not if_match:
            # Blind overwrite of an existing resource: refused on purpose.
            return (428, None)
        if if_match != entry["etag"]:
            # Stale precondition -- state is left exactly as it was.
            return (412, None)

        etag = self._next_etag()
        entry["etag"] = etag
        entry["body"] = body
        return (200, etag)

    def delete(self, key):
        self.items.pop(key, None)
