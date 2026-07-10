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

// RightSideView returns the value of the rightmost node at each
// depth of root's tree, top to bottom.
std::vector<int> RightSideView(TreeNode* root) {
    return {};
}
