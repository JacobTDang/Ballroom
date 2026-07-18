package main

// Update applies mask's dotted paths, copying values from source into
// target in place. Values are map[string]any where each entry is
// either a nested map[string]any (an object) or a string (a leaf).
//
// TODO: ignores the mask entirely -- just shallow-merges source into
// target, top-level keys only. No clearing, no path validation, no
// recursive per-path merge. Every rule in the problem statement is
// still yours to build.
func Update(target, source map[string]any, mask []string) error {
	for k, v := range source {
		target[k] = v
	}
	return nil
}
