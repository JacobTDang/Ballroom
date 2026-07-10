#include <climits>
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

// GoodNodes counts nodes X in root's tree where no node on the path
// from root to X has a value greater than X.
int GoodNodes(TreeNode* root) {
    std::function<int(TreeNode*, int)> dfs = [&](TreeNode* node, int maxSoFar) -> int {
        if (node == nullptr) return 0;
        int count = 0;
        if (node->val >= maxSoFar) {
            count = 1;
            maxSoFar = node->val;
        }
        count += dfs(node->left, maxSoFar);
        count += dfs(node->right, maxSoFar);
        return count;
    };
    return dfs(root, INT_MIN);
}
