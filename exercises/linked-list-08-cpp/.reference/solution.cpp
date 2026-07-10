#include <vector>

// FindDuplicate returns the one repeated value in nums, using Floyd's
// cycle detection over the implicit index -> nums[index] linked list.
int FindDuplicate(const std::vector<int>& nums) {
    int slow = nums[0], fast = nums[0];
    do {
        slow = nums[slow];
        fast = nums[nums[fast]];
    } while (slow != fast);

    int slow2 = nums[0];
    while (slow2 != slow) {
        slow2 = nums[slow2];
        slow = nums[slow];
    }
    return slow;
}
