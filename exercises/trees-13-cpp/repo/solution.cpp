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

// BuildTree reconstructs the unique binary tree whose preorder and
// inorder traversals are preorder and inorder.
TreeNode* BuildTree(std::vector<int>& preorder, std::vector<int>& inorder) {
    return nullptr;
}
