package main

import (
	"fmt"
	"strings"
)

// sourceLookup walks segments through source. found=false (absence
// anywhere along the walk) means the mask entry clears rather than
// sets -- that's a legitimate outcome, not an error.
func sourceLookup(source map[string]any, segments []string) (found bool, value any) {
	var node any = source
	for _, seg := range segments {
		m, ok := node.(map[string]any)
		if !ok {
			return false, nil
		}
		v, ok := m[seg]
		if !ok {
			return false, nil
		}
		node = v
	}
	return true, node
}

// targetParent walks every segment except the last through target --
// each one must already exist as an object. Returns the map the
// final segment lives in.
func targetParent(target map[string]any, segments []string, fullPath string) (map[string]any, error) {
	node := target
	for _, seg := range segments[:len(segments)-1] {
		v, ok := node[seg]
		if !ok {
			return nil, fmt.Errorf("unknown path %q: %q does not exist", fullPath, seg)
		}
		m, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unknown path %q: %q is not an object", fullPath, seg)
		}
		node = m
	}
	return node, nil
}

type maskOp struct {
	parent map[string]any
	leaf   string
	found  bool
	value  any
}

// Update: two passes on purpose. Validate every path's target-side
// intermediates first, THEN apply. A bad path anywhere in the mask
// must leave target completely untouched, not partially patched.
func Update(target, source map[string]any, mask []string) error {
	if len(mask) == 0 {
		return fmt.Errorf("update_mask must not be empty")
	}

	ops := make([]maskOp, 0, len(mask))
	for _, path := range mask {
		segments := strings.Split(path, ".")
		parent, err := targetParent(target, segments, path)
		if err != nil {
			return err
		}
		leaf := segments[len(segments)-1]
		found, value := sourceLookup(source, segments)
		ops = append(ops, maskOp{parent: parent, leaf: leaf, found: found, value: value})
	}

	for _, o := range ops {
		if o.found {
			o.parent[o.leaf] = o.value
		} else {
			delete(o.parent, o.leaf)
		}
	}
	return nil
}
