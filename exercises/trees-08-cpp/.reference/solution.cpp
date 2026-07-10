#include <vector>

// TreeNode is a binary tree node.
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

// LevelOrder returns root's node values grouped by depth, level by
// level from top to bottom, left to right within each level.
std::vector<std::vector<int>> LevelOrder(TreeNode* root) {
    std::vector<std::vector<int>> res;
    if (root == nullptr) return res;
    std::vector<TreeNode*> queue = {root};
    while (!queue.empty()) {
        std::vector<TreeNode*> next;
        std::vector<int> level;
        for (TreeNode* node : queue) {
            level.push_back(node->val);
            if (node->left != nullptr) next.push_back(node->left);
            if (node->right != nullptr) next.push_back(node->right);
        }
        res.push_back(level);
        queue = next;
    }
    return res;
}
