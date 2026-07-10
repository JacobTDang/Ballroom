#include <functional>
#include <vector>

// CanFinish reports whether all numCourses courses can be completed
// given the prerequisite pairs, i.e. whether the prerequisite graph
// has no cycle.
bool CanFinish(int numCourses, std::vector<std::vector<int>>& prerequisites) {
    std::vector<std::vector<int>> adj(numCourses);
    for (auto& p : prerequisites) {
        adj[p[0]].push_back(p[1]);
    }

    const int kUnvisited = 0, kVisiting = 1, kVisited = 2;
    std::vector<int> state(numCourses, kUnvisited);

    std::function<bool(int)> dfs = [&](int course) -> bool {
        if (state[course] == kVisiting) return false;
        if (state[course] == kVisited) return true;
        state[course] = kVisiting;
        for (int pre : adj[course]) {
            if (!dfs(pre)) return false;
        }
        state[course] = kVisited;
        return true;
    };

    for (int c = 0; c < numCourses; c++) {
        if (!dfs(c)) return false;
    }
    return true;
}
