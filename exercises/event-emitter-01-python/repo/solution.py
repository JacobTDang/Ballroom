class Emitter:
    """On/once subscribe (returning an id), off unsubscribes, emit
    calls the event's handlers in registration order.

    TODO: no ids (always 0), off does nothing, and once is just on --
    it never unhooks itself.
    """

    def __init__(self):
        self.handlers = {}

    def on(self, event, fn):
        self.handlers.setdefault(event, []).append(fn)
        return 0

    def once(self, event, fn):
        return self.on(event, fn)

    def off(self, id):
        pass

    def emit(self, event, value):
        for fn in self.handlers.get(event, []):
            fn(value)
