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

// KthSmallest returns the kth smallest value (1-indexed) among all
// nodes in the BST rooted at root.
int KthSmallest(TreeNode* root, int k) {
    std::vector<TreeNode*> stack;
    TreeNode* cur = root;
    while (true) {
        while (cur != nullptr) {
            stack.push_back(cur);
            cur = cur->left;
        }
        cur = stack.back();
        stack.pop_back();
        k--;
        if (k == 0) return cur->val;
        cur = cur->right;
    }
}
