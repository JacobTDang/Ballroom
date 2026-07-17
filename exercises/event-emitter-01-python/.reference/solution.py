class Emitter:
    """Every subscription gets a unique id in a single ordered list
    per event. emit walks a snapshot of ids and re-checks liveness
    before each call -- that's what makes removal-during-emit safe.
    once unhooks itself BEFORE calling, so re-entrant emits can't
    double-fire it."""

    def __init__(self):
        self.subs = {}  # event -> list of [id, fn, once]
        self.by_id = {}  # id -> event
        self.next_id = 0

    def _add(self, event, fn, once):
        self.next_id += 1
        self.subs.setdefault(event, []).append([self.next_id, fn, once])
        self.by_id[self.next_id] = event
        return self.next_id

    def on(self, event, fn):
        return self._add(event, fn, False)

    def once(self, event, fn):
        return self._add(event, fn, True)

    def off(self, id):
        event = self.by_id.pop(id, None)
        if event is None:
            return
        self.subs[event] = [s for s in self.subs[event] if s[0] != id]

    def emit(self, event, value):
        ids = [s[0] for s in self.subs.get(event, [])]
        for sid in ids:
            if self.by_id.get(sid) != event:
                continue  # removed during this emit
            sub = next((s for s in self.subs[event] if s[0] == sid), None)
            if sub is None:
                continue
            if sub[2]:
                self.off(sid)
            sub[1](value)
