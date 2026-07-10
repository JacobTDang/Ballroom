#include <functional>
#include <string>
#include <unordered_map>
#include <vector>

// trieNode is a prefix-tree node used to search for every word
// simultaneously during the board DFS.
struct trieNode {
    std::unordered_map<char, trieNode*> children;
    std::string word;  // non-empty exactly at nodes completing a full word
};

// FindWords returns every word from words that can be traced out on
// board via sequentially adjacent cells, each cell used at most once
// per word.
std::vector<std::string> FindWords(std::vector<std::vector<char>>& board,
                                    std::vector<std::string>& words) {
    trieNode* root = new trieNode();
    for (const auto& w : words) {
        trieNode* node = root;
        for (char c : w) {
            if (node->children.find(c) == node->children.end()) {
                node->children[c] = new trieNode();
            }
            node = node->children[c];
        }
        node->word = w;
    }

    int rows = static_cast<int>(board.size());
    int cols = static_cast<int>(board[0].size());
    std::vector<std::string> res;

    std::function<void(int, int, trieNode*)> dfs = [&](int r, int c, trieNode* node) {
        if (r < 0 || r >= rows || c < 0 || c >= cols) return;
        char ch = board[r][c];
        if (ch == '#') return;
        auto it = node->children.find(ch);
        if (it == node->children.end()) return;
        trieNode* next = it->second;
        if (!next->word.empty()) {
            res.push_back(next->word);
            next->word.clear();
        }
        board[r][c] = '#';
        dfs(r + 1, c, next);
        dfs(r - 1, c, next);
        dfs(r, c + 1, next);
        dfs(r, c - 1, next);
        board[r][c] = ch;
    };

    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            dfs(r, c, root);
        }
    }
    return res;
}
