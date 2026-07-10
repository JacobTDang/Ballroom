#include <cstdlib>
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

// IsBalanced reports whether every node's left and right subtrees
// differ in height by no more than 1.
bool IsBalanced(TreeNode* root) {
    // height returns -1 as a sentinel meaning "already found an
    // imbalance somewhere below", short-circuiting the rest of the walk.
    std::function<int(TreeNode*)> height = [&](TreeNode* n) -> int {
        if (n == nullptr) return 0;
        int l = height(n->left);
        if (l == -1) return -1;
        int r = height(n->right);
        if (r == -1) return -1;
        if (std::abs(l - r) > 1) return -1;
        return 1 + (l > r ? l : r);
    };
    return height(root) != -1;
}
