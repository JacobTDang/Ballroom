import base64


def _checksum(payload: str) -> str:
    h = 0
    for ch in payload:
        h = (h * 131 + ord(ch)) % 1_000_000_007
    return format(h, "x")


def _encode_token(last_id: int, page_size: int) -> str:
    payload = f"{last_id}:{page_size}"
    raw = f"{payload}:{_checksum(payload)}"
    return base64.urlsafe_b64encode(raw.encode()).decode()


def _decode_token(token: str):
    try:
        raw = base64.urlsafe_b64decode(token.encode()).decode()
    except Exception as exc:
        raise ValueError(f"invalid page token: {token!r}") from exc

    parts = raw.split(":")
    if len(parts) != 3:
        raise ValueError(f"invalid page token: {token!r}")

    payload = f"{parts[0]}:{parts[1]}"
    if _checksum(payload) != parts[2]:
        raise ValueError(f"invalid page token: {token!r} (checksum mismatch)")

    try:
        return int(parts[0]), int(parts[1])
    except ValueError as exc:
        raise ValueError(f"invalid page token: {token!r}") from exc


class CursorStore:
    """Keyset pagination: the token names the last id seen (plus the
    page_size it was issued for), never a raw offset -- so a walk
    can't be thrown off by inserts, and resuming with a different
    page_size is a loud error rather than a silent behavior change."""

    DEFAULT_PAGE_SIZE = 10
    MAX_PAGE_SIZE = 50

    def __init__(self, records):
        self.records = {r["id"]: dict(r) for r in records}

    def insert(self, record):
        if record["id"] in self.records:
            raise ValueError(f"duplicate id: {record['id']}")
        self.records[record["id"]] = dict(record)

    def list(self, page_size, page_token=""):
        if page_size is None or page_size <= 0:
            effective = self.DEFAULT_PAGE_SIZE
        elif page_size > self.MAX_PAGE_SIZE:
            effective = self.MAX_PAGE_SIZE
        else:
            effective = page_size

        if page_token:
            cursor_id, encoded_size = _decode_token(page_token)
            if encoded_size != effective:
                raise ValueError(
                    f"page_token was issued for page_size={encoded_size}, "
                    f"not {effective}"
                )
        else:
            cursor_id = None

        candidates = sorted(
            (r for r in self.records.values() if cursor_id is None or r["id"] > cursor_id),
            key=lambda r: r["id"],
        )
        page = candidates[:effective]
        if len(candidates) > effective:
            next_token = _encode_token(page[-1]["id"], effective)
        else:
            next_token = ""
        return page, next_token
