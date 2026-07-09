# LRU Cache

Design a data structure that follows the constraints of a **Least
Recently Used (LRU) cache**.

Implement `LRUCache`:

- `LRUCache(capacity)` initializes the cache with positive size `capacity`.
- `get(key)` returns the value of `key` if it exists, otherwise returns
  "not found" (`-1` in C++/Python, `(0, false)` in Go). Accessing a key
  counts as a "use" of it.
- `put(key, value)` updates the value of `key` if it exists, otherwise
  adds the key-value pair. If adding a new key would exceed `capacity`,
  evict the **least recently used** key first.

Both `get` and `put` must run in **O(1)** average time.

## Example

```
cache = LRUCache(2)
cache.put(1, 1)
cache.put(2, 2)
cache.get(1)     // returns 1, and marks 1 as most recently used
cache.put(3, 3)  // evicts key 2 (least recently used)
cache.get(2)     // returns "not found"
cache.put(4, 4)  // evicts key 1
cache.get(1)     // returns "not found"
cache.get(3)     // returns 3
cache.get(4)     // returns 4
```

## Constraints

- `1 <= capacity <= 3000`
- `0 <= key, value <= 10^4`
- At most `2 * 10^5` calls will be made to `get` and `put`.
