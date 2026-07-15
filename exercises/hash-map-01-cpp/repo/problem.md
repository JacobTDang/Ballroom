# Design a Hash Map

Design a hash map from scratch — without using your language's
built-in hash table (no `dict`/`map`/`unordered_map` for the storage;
arrays/lists/vectors are fine).

Support non-negative integer keys and values:

- `put(key, value)` — insert, or update if the key exists
- `get(key)` — return the value, or `-1` if the key is absent
- `remove(key)` — delete the key if present

## Examples

```
put(1, 100)
get(1)      -> 100
put(1, 200)
get(1)      -> 200
remove(1)
get(1)      -> -1
```

## Constraints

- `0 <= key, value <= 10^6`
- Design for many keys: think buckets and collision handling, not one
  giant array.
