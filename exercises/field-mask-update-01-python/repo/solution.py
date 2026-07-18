def update(target, source, mask):
    """Apply mask's dotted paths, copying values from source into
    target in place.

    TODO: ignores the mask entirely -- just shallow-merges source into
    target, top-level keys only. No clearing, no path validation, no
    recursive per-path merge. Every rule in the problem statement is
    still yours to build.
    """
    target.update(source)
