#include <queue>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>

// LadderLength returns the number of words in the shortest
// transformation sequence from beginWord to endWord, changing one
// letter at a time through words in wordList, or 0 if impossible.
int LadderLength(std::string beginWord, std::string endWord, std::vector<std::string>& wordList) {
    std::unordered_set<std::string> wordSet(wordList.begin(), wordList.end());
    if (!wordSet.count(endWord)) return 0;

    std::unordered_map<std::string, std::vector<std::string>> patterns;
    auto addPatterns = [&](const std::string& word) {
        for (size_t i = 0; i < word.size(); i++) {
            std::string pattern = word;
            pattern[i] = '*';
            patterns[pattern].push_back(word);
        }
    };
    for (const auto& w : wordSet) addPatterns(w);
    addPatterns(beginWord);

    std::unordered_set<std::string> visited{beginWord};
    std::queue<std::string> queue;
    queue.push(beginWord);
    int steps = 1;

    while (!queue.empty()) {
        int levelSize = static_cast<int>(queue.size());
        for (int i = 0; i < levelSize; i++) {
            std::string word = queue.front();
            queue.pop();
            if (word == endWord) return steps;
            for (size_t j = 0; j < word.size(); j++) {
                std::string pattern = word;
                pattern[j] = '*';
                for (const auto& neighbor : patterns[pattern]) {
                    if (!visited.count(neighbor)) {
                        visited.insert(neighbor);
                        queue.push(neighbor);
                    }
                }
            }
        }
        steps++;
    }
    return 0;
}
