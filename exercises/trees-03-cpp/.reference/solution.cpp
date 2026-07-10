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

// DiameterOfBinaryTree returns the number of edges on the longest
// path between any two nodes in root's tree.
int DiameterOfBinaryTree(TreeNode* root) {
    int best = 0;
    std::function<int(TreeNode*)> height = [&](TreeNode* n) -> int {
        if (n == nullptr) return 0;
        int l = height(n->left);
        int r = height(n->right);
        best = std::max(best, l + r);
        return 1 + std::max(l, r);
    };
    height(root);
    return best;
}
