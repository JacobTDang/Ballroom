#include <functional>
#include <sstream>
#include <string>
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

// Serialize encodes root as a string that Deserialize can turn back
// into an equivalent tree. The exact format is up to you.
std::string Serialize(TreeNode* root) {
    std::string out;
    std::function<void(TreeNode*)> walk = [&](TreeNode* node) {
        if (node == nullptr) {
            out += "#,";
            return;
        }
        out += std::to_string(node->val) + ",";
        walk(node->left);
        walk(node->right);
    };
    walk(root);
    return out;
}

// Deserialize decodes a string produced by Serialize back into the
// original tree.
TreeNode* Deserialize(const std::string& data) {
    std::vector<std::string> tokens;
    std::stringstream ss(data);
    std::string tok;
    while (std::getline(ss, tok, ',')) {
        tokens.push_back(tok);
    }

    size_t idx = 0;
    std::function<TreeNode*()> walk = [&]() -> TreeNode* {
        if (idx >= tokens.size() || tokens[idx] == "#") {
            idx++;
            return nullptr;
        }
        int v = std::stoi(tokens[idx]);
        idx++;
        TreeNode* node = new TreeNode(v);
        node->left = walk();
        node->right = walk();
        return node;
    };
    return walk();
}
