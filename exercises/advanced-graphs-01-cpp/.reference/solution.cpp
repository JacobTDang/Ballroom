#include <algorithm>
#include <functional>
#include <map>
#include <string>
#include <vector>

// FindItinerary reconstructs an itinerary from tickets, all starting from
// "JFK", using every ticket exactly once, returning the lexicographically
// smallest valid sequence of airports.
std::vector<std::string> FindItinerary(std::vector<std::vector<std::string>>& tickets) {
    std::map<std::string, std::vector<std::string>> graph;
    for (auto& t : tickets) {
        graph[t[0]].push_back(t[1]);
    }
    for (auto& entry : graph) {
        std::sort(entry.second.begin(), entry.second.end());
    }

    std::vector<std::string> route;
    std::function<void(const std::string&)> visit = [&](const std::string& airport) {
        auto& dests = graph[airport];
        while (!dests.empty()) {
            std::string next = dests.front();
            dests.erase(dests.begin());
            visit(next);
        }
        route.push_back(airport);
    };
    visit("JFK");

    std::reverse(route.begin(), route.end());
    return route;
}
