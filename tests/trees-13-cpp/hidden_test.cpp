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

TreeNode* BuildTree(std::vector<int>& preorder, std::vector<int>& inorder);

using OptInt = std::optional<int>;

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

void check(std::vector<int> preorder, std::vector<int> inorder, std::vector<OptInt> want) {
    assert(toLevelOrder(BuildTree(preorder, inorder)) == want);
}

int main() {
    check({3, 9, 20, 15, 7}, {9, 3, 15, 20, 7}, {3, 9, 20, std::nullopt, std::nullopt, 15, 7});
    check({-1}, {-1}, {-1});
    check({1, 2, 3}, {3, 2, 1}, {1, 2, std::nullopt, 3});
    check({1, 2, 4, 5, 3, 6, 7}, {4, 2, 5, 1, 6, 3, 7}, {1, 2, 3, 4, 5, 6, 7});
    check({1, 2}, {1, 2}, {1, std::nullopt, 2});
    printf("all assertions passed\n");
    return 0;
}
