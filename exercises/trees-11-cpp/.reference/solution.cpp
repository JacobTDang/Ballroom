#include <cstdint>
#include <functional>
#include <limits>

// TreeNode is a binary tree node.
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

// IsValidBST reports whether root is a valid binary search tree.
bool IsValidBST(TreeNode* root) {
    std::function<bool(TreeNode*, int64_t, int64_t)> valid = [&](TreeNode* node, int64_t lo,
                                                                   int64_t hi) -> bool {
        if (node == nullptr) return true;
        int64_t v = node->val;
        if (v <= lo || v >= hi) return false;
        return valid(node->left, lo, v) && valid(node->right, v, hi);
    };
    return valid(root, std::numeric_limits<int64_t>::min(), std::numeric_limits<int64_t>::max());
}
