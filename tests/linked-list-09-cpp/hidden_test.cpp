#include <cassert>
#include <cstdio>
#include <vector>

struct ListNode {
    int val;
    ListNode* next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode* next) : val(x), next(next) {}
};

ListNode* MergeKLists(std::vector<ListNode*>& lists);

ListNode* buildList(const std::vector<int>& vals) {
    ListNode dummy;
    ListNode* cur = &dummy;
    for (int v : vals) {
        cur->next = new ListNode(v);
        cur = cur->next;
    }
    return dummy.next;
}

std::vector<int> toVector(ListNode* head) {
    std::vector<int> out;
    for (ListNode* n = head; n != nullptr; n = n->next) {
        out.push_back(n->val);
    }
    return out;
}

void check(std::vector<std::vector<int>> in, std::vector<int> want) {
    std::vector<ListNode*> lists;
    for (auto& v : in) lists.push_back(buildList(v));
    assert(toVector(MergeKLists(lists)) == want);
}

int main() {
    check({{1, 4, 5}, {1, 3, 4}, {2, 6}}, {1, 1, 2, 3, 4, 4, 5, 6});
    check({}, {});
    check({{}}, {});
    check({{1}, {}, {2}}, {1, 2});
    printf("all assertions passed\n");
    return 0;
}
