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
    return 0;
}
