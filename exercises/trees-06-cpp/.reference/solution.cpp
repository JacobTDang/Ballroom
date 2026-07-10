// TreeNode is a binary tree node.
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

bool isSameTree(TreeNode* p, TreeNode* q) {
    if (p == nullptr && q == nullptr) return true;
    if (p == nullptr || q == nullptr || p->val != q->val) return false;
    return isSameTree(p->left, q->left) && isSameTree(p->right, q->right);
}

// IsSubtree reports whether subRoot's tree matches some node in
// root's tree and everything below it.
bool IsSubtree(TreeNode* root, TreeNode* subRoot) {
    if (root == nullptr) return subRoot == nullptr;
    if (isSameTree(root, subRoot)) return true;
    return IsSubtree(root->left, subRoot) || IsSubtree(root->right, subRoot);
}
