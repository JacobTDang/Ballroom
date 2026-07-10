#include <functional>
#include <unordered_map>
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
    std::unordered_map<int, int> inorderIdx;
    for (size_t i = 0; i < inorder.size(); i++) inorderIdx[inorder[i]] = static_cast<int>(i);

    int pre = 0;
    std::function<TreeNode*(int, int)> build = [&](int inLo, int inHi) -> TreeNode* {
        if (inLo > inHi) return nullptr;
        int rootVal = preorder[pre++];
        TreeNode* root = new TreeNode(rootVal);
        int mid = inorderIdx[rootVal];
        root->left = build(inLo, mid - 1);
        root->right = build(mid + 1, inHi);
        return root;
    };
    return build(0, static_cast<int>(inorder.size()) - 1);
}
