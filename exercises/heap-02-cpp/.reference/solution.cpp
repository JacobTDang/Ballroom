#include <queue>
#include <vector>

// LastStoneWeight repeatedly smashes the two heaviest stones
// together and returns the weight of whatever stone (if any) remains.
int LastStoneWeight(std::vector<int>& stones) {
    std::priority_queue<int> heap(stones.begin(), stones.end());
    while (heap.size() > 1) {
        int a = heap.top();
        heap.pop();
        int b = heap.top();
        heap.pop();
        if (a != b) heap.push(a - b);
    }
    return heap.empty() ? 0 : heap.top();
}
