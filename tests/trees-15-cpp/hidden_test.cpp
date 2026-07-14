#include <cassert>
#include <cstdio>
#include <optional>
#include <string>
#include <vector>

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

std::string Serialize(TreeNode* root);
TreeNode* Deserialize(const std::string& data);

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

std::vector<OptInt> toLevelOrder(TreeNode* root) {
    if (root == nullptr) return {};
    std::vector<OptInt> out;
    std::vector<TreeNode*> queue = {root};
    while (!queue.empty()) {
        TreeNode* node = queue.front();
        queue.erase(queue.begin());
        if (node == nullptr) {
            out.push_back(std::nullopt);
            continue;
        }
        out.push_back(node->val);
        queue.push_back(node->left);
        queue.push_back(node->right);
    }
    while (!out.empty() && !out.back().has_value()) out.pop_back();
    return out;
}

void check(std::vector<OptInt> vals) {
    TreeNode* original = buildTree(vals);
    TreeNode* roundTripped = Deserialize(Serialize(original));
    assert(toLevelOrder(roundTripped) == toLevelOrder(original));
}

int main() {
    check({1, 2, 3, std::nullopt, std::nullopt, 4, 5});
    check({});
    check({1});
    check({-1, -2, -3});
    check({5, 4, 7, 3, std::nullopt, 2, std::nullopt, -1, std::nullopt, 9});
    check({0});
    check({100, std::nullopt, 200, std::nullopt, std::nullopt, std::nullopt, 300});
    printf("all assertions passed\n");
    return 0;
}
