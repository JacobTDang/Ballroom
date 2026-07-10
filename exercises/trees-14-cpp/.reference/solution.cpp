#include <algorithm>
#include <functional>

// TreeNode is a binary tree node.
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

// MaxPathSum returns the maximum sum of any non-empty path between
// two nodes in root's tree (the path need not pass through root).
int MaxPathSum(TreeNode* root) {
    int best = root->val;
    std::function<int(TreeNode*)> gain = [&](TreeNode* node) -> int {
        if (node == nullptr) return 0;
        int leftGain = std::max(gain(node->left), 0);
        int rightGain = std::max(gain(node->right), 0);
        best = std::max(best, node->val + leftGain + rightGain);
        return node->val + std::max(leftGain, rightGain);
    };
    gain(root);
    return best;
}
