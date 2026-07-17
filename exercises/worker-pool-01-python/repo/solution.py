def process_all(jobs, workers, fn):
    """Apply fn to every job, returning results in input order, using
    at most `workers` concurrent worker threads.

    TODO: this version is sequential -- one job at a time, no threads
    at all. Parallelize it without breaking the ordering.
    """
    return [fn(j) for j in jobs]
