#include <cassert>
#include <cstdio>
#include <optional>
#include <vector>

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

TreeNode* LowestCommonAncestor(TreeNode* root, TreeNode* p, TreeNode* q);

using OptInt = std::optional<int>;

TreeNode* buildTree(const std::vector<OptInt>& vals) {
    if (vals.empty() || !vals[0].has_value()) return nullptr;
    TreeNode* root = new TreeNode(*vals[0]);
    std::vector<TreeNode*> queue = {root};
    size_t i = 1;
    while (!queue.empty() && i < vals.size()) {
        TreeNode* node = queue.front();
        queue.erase(queue.begin());
        if (i < vals.size()) {
            if (vals[i].has_value()) {
                node->left = new TreeNode(*vals[i]);
                queue.push_back(node->left);
            }
            i++;
        }
        if (i < vals.size()) {
            if (vals[i].has_value()) {
                node->right = new TreeNode(*vals[i]);
                queue.push_back(node->right);
            }
            i++;
        }
    }
    return root;
}

TreeNode* findNode(TreeNode* root, int val) {
    while (root != nullptr) {
        if (val == root->val) return root;
        root = (val < root->val) ? root->left : root->right;
    }
    return nullptr;
}

int main() {
    std::vector<OptInt> tree = {6, 2, 8, 0, 4, 7, 9, std::nullopt, std::nullopt, 3, 5};

    {
        TreeNode* root = buildTree(tree);
        assert(LowestCommonAncestor(root, findNode(root, 2), findNode(root, 8))->val == 6);
    }
    {
        TreeNode* root = buildTree(tree);
        assert(LowestCommonAncestor(root, findNode(root, 2), findNode(root, 4))->val == 2);
    }
    {
        TreeNode* root = buildTree(tree);
        assert(LowestCommonAncestor(root, findNode(root, 0), findNode(root, 5))->val == 2);
    }
    {
        TreeNode* root = buildTree(tree);
        assert(LowestCommonAncestor(root, findNode(root, 7), findNode(root, 9))->val == 8);
    }
    {
        TreeNode* root = buildTree(tree);
        assert(LowestCommonAncestor(root, findNode(root, 6), findNode(root, 6))->val == 6);
    }
    {
        TreeNode* root = buildTree(tree);
        assert(LowestCommonAncestor(root, findNode(root, 0), findNode(root, 3))->val == 2);
    }
    {
        TreeNode* root = buildTree({2, 1});
        assert(LowestCommonAncestor(root, findNode(root, 2), findNode(root, 1))->val == 2);
    }
    printf("all assertions passed\n");
    return 0;
}
