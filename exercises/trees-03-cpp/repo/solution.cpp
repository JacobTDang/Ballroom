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
    return 0;
}
