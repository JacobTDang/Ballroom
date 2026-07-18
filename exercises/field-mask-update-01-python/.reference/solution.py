def _source_lookup(source, segments):
    """Walk segments through source. found=False (absence anywhere
    along the walk) means the mask entry clears rather than sets --
    that's a legitimate outcome, not an error."""
    node = source
    for seg in segments:
        if not isinstance(node, dict) or seg not in node:
            return False, None
        node = node[seg]
    return True, node


def _target_parent(target, segments, full_path):
    """Walk every segment except the last through target -- each one
    must already exist as an object. Raises ValueError naming
    full_path if not. Returns the dict the final segment lives in."""
    node = target
    for seg in segments[:-1]:
        if not isinstance(node, dict) or seg not in node:
            raise ValueError(f"unknown path {full_path!r}: {seg!r} does not exist")
        node = node[seg]
        if not isinstance(node, dict):
            raise ValueError(f"unknown path {full_path!r}: {seg!r} is not an object")
    return node


def update(target, source, mask):
    """Two passes on purpose: validate every path's target-side
    intermediates first, THEN apply. A bad path anywhere in the mask
    must leave target completely untouched, not partially patched."""
    if not mask:
        raise ValueError("update_mask must not be empty")

    ops = []  # (parent_dict, leaf_key, found_in_source, value)
    for path in mask:
        segments = path.split(".")
        parent = _target_parent(target, segments, path)
        leaf = segments[-1]
        found, value = _source_lookup(source, segments)
        ops.append((parent, leaf, found, value))

    for parent, leaf, found, value in ops:
        if found:
            parent[leaf] = value
        else:
            parent.pop(leaf, None)
