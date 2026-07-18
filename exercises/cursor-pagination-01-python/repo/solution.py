class CursorStore:
    """Keyset pagination over id-sorted records.

    TODO: this paginates by raw position offset -- it looks fine until
    the data changes between calls. No tamper check, no page-size
    clamp, no page-size binding in the token. Every rule in the
    problem statement is still yours to build.
    """

    DEFAULT_PAGE_SIZE = 10
    MAX_PAGE_SIZE = 50

    def __init__(self, records):
        self.records = {r["id"]: dict(r) for r in records}

    def insert(self, record):
        self.records[record["id"]] = dict(record)

    def list(self, page_size, page_token=""):
        offset = int(page_token) if page_token else 0
        ordered = sorted(self.records.values(), key=lambda r: r["id"])
        page = ordered[offset:offset + page_size]
        next_token = str(offset + page_size) if offset + page_size < len(ordered) else ""
        return page, next_token
