// TreeNode is a binary tree node.
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

// GoodNodes counts nodes X in root's tree where no node on the path
// from root to X has a value greater than X.
int GoodNodes(TreeNode* root) {
    return 0;
}
