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

bool IsSameTree(TreeNode* p, TreeNode* q);

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

int main() {
    assert(IsSameTree(buildTree({1, 2, 3}), buildTree({1, 2, 3})) == true);
    assert(IsSameTree(buildTree({1, 2}), buildTree({1, std::nullopt, 2})) == false);
    assert(IsSameTree(buildTree({1, 2, 1}), buildTree({1, 1, 2})) == false);
    assert(IsSameTree(buildTree({}), buildTree({})) == true);
    assert(IsSameTree(buildTree({1}), buildTree({})) == false);
    printf("all assertions passed\n");
    return 0;
}
