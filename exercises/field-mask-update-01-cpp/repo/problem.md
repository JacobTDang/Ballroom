# Field Mask Update Engine

AIP-134's `update_mask`: a PATCH names *exactly* which fields change by
listing their dotted paths, and nothing else moves — not sibling
fields, not fields `source` happens to also carry but that aren't in
the mask.

`Update(target, source, mask)` applies `mask` (a list of dotted paths
like `"address.city"`) by copying each path's value from `source` into
`target`, in place. The starter ships a `Value` type (keep it as your
value representation — building yet another JSON-ish variant isn't the
exercise): a `Value` is either an **object** (`std::map<std::string,
Value>`) or a **scalar** (`std::string`).

The rules:

- **A path in the mask that's missing from `source`** (the leaf itself
  absent, or a whole ancestor object absent) means *clear that field*:
  delete the key from `target`. This is how a field mask sets
  something back to empty, not by copying an explicit empty value.
- **Every path segment except the last must already exist in `target`
  as an object.** The mask updates an existing resource's fields — it
  doesn't invent new nested structure. Only the *last* segment (the
  leaf actually being written or cleared) may be new.
- **An intermediate segment that's missing, or that exists but isn't
  an object** (you can't descend into a leaf value), is an unknown
  path: a loud error that names the exact path that failed.
- **An empty mask is an error** — "update nothing" is almost always a
  caller mistake, not a valid request.
- Paths are literal dotted field names — `*` is not a wildcard here,
  just an ordinary (almost certainly unknown) path segment.
- Nothing outside the masked paths is touched: sibling fields at every
  level are left exactly as they were, and a path failing validation
  leaves `target` completely unchanged (validate every path in the
  mask before applying any of them).

## The invariant the tests enforce

Updating one leaf never disturbs its siblings at any level; several
mask paths (including ones sharing a parent) apply together correctly;
an omitted path clears rather than leaving the old value; a missing
intermediate and a scalar-typed intermediate are both loud errors
naming the offending path; an empty mask is a loud error.

API: `struct Value { enum Kind { OBJECT, SCALAR } kind; std::map<std::string, Value> object; std::string scalar; }` (provided — see `IsObject()`/`Value::Obj(...)`/`Value::Str(...)` helpers in the starter); `bool Update(Value* target, const Value& source, const std::vector<std::string>& mask, std::string* err)` — mutates `*target` in place; returns `false` and fills `*err` (naming the offending path) on an empty mask or an unresolvable intermediate segment.
