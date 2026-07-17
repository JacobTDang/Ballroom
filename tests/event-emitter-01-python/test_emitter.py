from solution import Emitter


def test_registration_order_and_event_isolation():
    e = Emitter()
    got = []
    e.on("a", lambda v: got.append("first"))
    e.on("a", lambda v: got.append("second"))
    e.on("b", lambda v: got.append("other-event"))
    e.emit("a", 1)
    assert got == ["first", "second"]


def test_once_fires_exactly_once():
    e = Emitter()
    calls = []
    e.once("a", lambda v: calls.append(v))
    e.emit("a", 1)
    e.emit("a", 2)
    assert len(calls) == 1


def test_off_removes_exactly_that_subscription():
    e = Emitter()
    got = []
    id = e.on("a", lambda v: got.append("removed"))
    e.on("a", lambda v: got.append("kept"))
    e.off(id)
    e.emit("a", 1)
    assert got == ["kept"]


def test_removal_during_emit_skips_the_removed():
    e = Emitter()
    got = []
    victim = []

    def assassin(v):
        got.append("assassin")
        e.off(victim[0])

    e.on("a", assassin)
    victim.append(e.on("a", lambda v: got.append("victim")))
    e.on("a", lambda v: got.append("bystander"))
    e.emit("a", 1)
    assert got == ["assassin", "bystander"], "the removed handler must not fire"


def test_emit_with_no_handlers_is_a_noop():
    Emitter().emit("nobody", 42)


def test_handlers_receive_the_value():
    e = Emitter()
    got = []
    e.on("a", got.append)
    e.emit("a", 99)
    assert got == [99]
