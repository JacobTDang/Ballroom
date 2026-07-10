#include <algorithm>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>

// AlienOrder derives a valid character ordering for the alien alphabet
// implied by words, assumed sorted according to that unknown ordering.
// Returns "" if no valid ordering exists.
std::string AlienOrder(std::vector<std::string>& words) {
    std::unordered_map<char, std::unordered_set<char>> adj;
    std::unordered_map<char, int> inDegree;

    for (auto& w : words) {
        for (char c : w) {
            if (adj.find(c) == adj.end()) {
                adj[c] = {};
                inDegree[c] = 0;
            }
        }
    }

    for (size_t i = 0; i + 1 < words.size(); i++) {
        std::string& w1 = words[i];
        std::string& w2 = words[i + 1];
        size_t minLen = std::min(w1.size(), w2.size());
        if (w1.size() > w2.size() && w1.substr(0, minLen) == w2.substr(0, minLen)) {
            return "";
        }
        for (size_t j = 0; j < minLen; j++) {
            if (w1[j] != w2[j]) {
                if (adj[w1[j]].find(w2[j]) == adj[w1[j]].end()) {
                    adj[w1[j]].insert(w2[j]);
                    inDegree[w2[j]]++;
                }
                break;
            }
        }
    }

    std::vector<char> queue;
    for (auto& entry : inDegree) {
        if (entry.second == 0) queue.push_back(entry.first);
    }
    std::sort(queue.begin(), queue.end());

    std::string order;
    size_t head = 0;
    while (head < queue.size()) {
        char c = queue[head++];
        order.push_back(c);

        std::vector<char> neighbors(adj[c].begin(), adj[c].end());
        std::sort(neighbors.begin(), neighbors.end());

        for (char n : neighbors) {
            inDegree[n]--;
            if (inDegree[n] == 0) {
                queue.push_back(n);
            }
        }
    }

    if (order.size() != inDegree.size()) {
        return "";
    }
    return order;
}
