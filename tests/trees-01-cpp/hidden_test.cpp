#include <cassert>
#include <cstdio>
#include <optional>
#include <vector>

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

TreeNode* InvertTree(TreeNode* root);

using OptInt = std::optional<int>;

// buildTree builds a binary tree from vals in LeetCode's level-order
// array format (nullopt entries are missing children).
TreeNode* buildTree(const std::vector<OptInt>& vals) {
    if (vals.empty() || !vals[0].has_value()) return nullptr;
    TreeNode* root = new TreeNode(*vals[0]);
    std::vector<TreeNode*> queue = {root};
    size_t i = 1;
    while (!queue.empty() && i < vals.size()) {
        TreeNode* node = queue.front();
        queue.erase(queue.begin());
        if (i < vals.size()) {
            if (vals[i].has_value()) {
                node->left = new TreeNode(*vals[i]);
                queue.push_back(node->left);
            }
            i++;
        }
        if (i < vals.size()) {
            if (vals[i].has_value()) {
                node->right = new TreeNode(*vals[i]);
                queue.push_back(node->right);
            }
            i++;
        }
    }
    return root;
}

// toLevelOrder serializes a tree back to the same nullopt-padded
// level-order format buildTree consumes, trimming only the trailing
// run of nullopts so results compare cleanly.
std::vector<OptInt> toLevelOrder(TreeNode* root) {
    if (root == nullptr) return {};
    std::vector<OptInt> out;
    std::vector<TreeNode*> queue = {root};
    while (!queue.empty()) {
        TreeNode* node = queue.front();
        queue.erase(queue.begin());
        if (node == nullptr) {
            out.push_back(std::nullopt);
            continue;
        }
        out.push_back(node->val);
        queue.push_back(node->left);
        queue.push_back(node->right);
    }
    while (!out.empty() && !out.back().has_value()) out.pop_back();
    return out;
}

void check(std::vector<OptInt> in, std::vector<OptInt> want) {
    assert(toLevelOrder(InvertTree(buildTree(in))) == want);
}

int main() {
    check({4, 2, 7, 1, 3, 6, 9}, {4, 7, 2, 9, 6, 3, 1});
    check({2, 1, 3}, {2, 3, 1});
    check({}, {});
    check({1}, {1});
    check({1, 2, std::nullopt}, {1, std::nullopt, 2});
    printf("all assertions passed\n");
    return 0;
}
