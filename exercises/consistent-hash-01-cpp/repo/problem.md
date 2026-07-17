# Consistent-Hash Ring

Map keys to nodes so that adding or removing one node only remaps the
keys in its neighborhood — the property that lets a cache cluster grow
without a stampede. Nodes sit at (many virtual) positions on a hash
ring; a key belongs to the first node clockwise from its hash.

The starter is `hash(key) % len(nodes)` — correct until the node count
changes, then almost every key moves. The tests measure exactly that.

## The invariant the tests enforce

- Deterministic: the same key always maps to the same node while the
  ring is unchanged.
- Balanced: with virtual nodes, 3 nodes each get a sane share of
  10,000 keys (no node under 10% or over 60%).
- **Minimal remap**: adding a 4th node moves some keys (more than 5%)
  but fewer than half — `%N` moves ~75% and fails. Removing that node
  restores the *exact* original mapping.
- An empty ring returns "" for every lookup.

API: `class Ring { Ring(int vnodes); void AddNode(const std::string&); void RemoveNode(const std::string&); std::string Lookup(const std::string&); }` ("" on an empty ring).

Think: what data structure answers "first node position clockwise of
this hash" quickly?
