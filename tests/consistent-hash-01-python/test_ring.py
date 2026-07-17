from solution import Ring


def build_ring():
    r = Ring(100)
    r.add_node("node-a")
    r.add_node("node-b")
    r.add_node("node-c")
    return r


def test_deterministic_lookup():
    r = build_ring()
    for i in range(100):
        key = f"key-{i}"
        first = r.lookup(key)
        assert first, "lookup returned no node on a populated ring"
        assert r.lookup(key) == first


def test_reasonable_balance():
    r = build_ring()
    counts = {}
    keys = 10000
    for i in range(keys):
        n = r.lookup(f"key-{i}")
        counts[n] = counts.get(n, 0) + 1
    for node in ("node-a", "node-b", "node-c"):
        share = counts.get(node, 0) / keys
        assert 0.10 <= share <= 0.60, f"{node} owns {share:.0%}, want 10-60% with 100 vnodes"


def test_add_remaps_neighborhood_and_remove_restores():
    r = build_ring()
    keys = 10000
    before = {f"key-{i}": r.lookup(f"key-{i}") for i in range(keys)}

    r.add_node("node-d")
    moved = sum(1 for k, owner in before.items() if r.lookup(k) != owner)
    assert moved * 2 < keys, \
        f"adding one node moved {moved}/{keys} keys -- %N rehashing moves nearly everything"
    assert moved * 20 >= keys, \
        f"adding a node moved only {moved}/{keys} keys -- it isn't taking its share"

    r.remove_node("node-d")
    for k, owner in before.items():
        assert r.lookup(k) == owner, "removing the node must restore the exact original mapping"


def test_empty_ring():
    assert Ring(100).lookup("anything") == ""
