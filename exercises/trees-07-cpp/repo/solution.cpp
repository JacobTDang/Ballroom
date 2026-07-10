// TreeNode is a binary tree node.
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

// LowestCommonAncestor returns the lowest node in the BST rooted at
// root that has both p and q as descendants (a node counts as its
// own descendant).
TreeNode* LowestCommonAncestor(TreeNode* root, TreeNode* p, TreeNode* q) {
    return nullptr;
}
