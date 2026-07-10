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
    return -1;
}
