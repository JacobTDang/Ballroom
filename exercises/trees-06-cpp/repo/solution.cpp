// TreeNode is a binary tree node.
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

// IsSubtree reports whether subRoot's tree matches some node in
// root's tree and everything below it.
bool IsSubtree(TreeNode* root, TreeNode* subRoot) {
    return false;
}
