#include <functional>
#include <vector>

// ProcessAll applies fn to every job and returns the results in input
// order, using at most `workers` concurrent workers.
//
// TODO: this version is sequential -- one job at a time, no threads at
// all. Parallelize it without breaking the ordering.
std::vector<int> ProcessAll(const std::vector<int>& jobs, int workers,
                            std::function<int(int)> fn) {
    std::vector<int> results(jobs.size());
    for (size_t i = 0; i < jobs.size(); i++) {
        results[i] = fn(jobs[i]);
    }
    return results;
}
