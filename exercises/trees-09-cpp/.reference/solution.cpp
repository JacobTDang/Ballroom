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

// RightSideView returns the value of the rightmost node at each
// depth of root's tree, top to bottom.
std::vector<int> RightSideView(TreeNode* root) {
    std::vector<int> res;
    if (root == nullptr) return res;
    std::vector<TreeNode*> queue = {root};
    while (!queue.empty()) {
        std::vector<TreeNode*> next;
        for (size_t i = 0; i < queue.size(); i++) {
            TreeNode* node = queue[i];
            if (i == queue.size() - 1) res.push_back(node->val);
            if (node->left != nullptr) next.push_back(node->left);
            if (node->right != nullptr) next.push_back(node->right);
        }
        queue = next;
    }
    return res;
}
