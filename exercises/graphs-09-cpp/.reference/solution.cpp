#include <functional>
#include <vector>

// FindOrder returns a valid course order satisfying every
// prerequisite pair, or an empty vector if no valid order exists.
std::vector<int> FindOrder(int numCourses, std::vector<std::vector<int>>& prerequisites) {
    std::vector<std::vector<int>> adj(numCourses);
    for (auto& p : prerequisites) {
        adj[p[0]].push_back(p[1]);
    }

    const int kUnvisited = 0, kVisiting = 1, kVisited = 2;
    std::vector<int> state(numCourses, kUnvisited);
    std::vector<int> order;
    order.reserve(numCourses);

    std::function<bool(int)> dfs = [&](int course) -> bool {
        if (state[course] == kVisiting) return false;
        if (state[course] == kVisited) return true;
        state[course] = kVisiting;
        for (int pre : adj[course]) {
            if (!dfs(pre)) return false;
        }
        state[course] = kVisited;
        order.push_back(course);
        return true;
    };

    for (int c = 0; c < numCourses; c++) {
        if (!dfs(c)) return {};
    }
    return order;
}
