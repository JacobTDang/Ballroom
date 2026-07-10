#include <string>

// TreeNode is a binary tree node.
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

// Serialize encodes root as a string that Deserialize can turn back
// into an equivalent tree. The exact format is up to you.
std::string Serialize(TreeNode* root) {
    return "";
}

// Deserialize decodes a string produced by Serialize back into the
// original tree.
TreeNode* Deserialize(const std::string& data) {
    return nullptr;
}
