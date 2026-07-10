#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<int> PartitionLabels(std::string s);

void testClassic() {
    std::vector<int> want = {9, 7, 8};
    assert(PartitionLabels("ababcbacadefegdehijhklij") == want);
}

void testAllUnique() {
    std::vector<int> want = {1, 1, 1, 1, 1};
    assert(PartitionLabels("abcde") == want);
}

void testAllSame() {
    std::vector<int> want = {4};
    assert(PartitionLabels("aaaa") == want);
}

void testSingleChar() {
    std::vector<int> want = {1};
    assert(PartitionLabels("a") == want);
}

void testMultipleEqualPartitions() {
    std::vector<int> want = {2, 2, 2};
    assert(PartitionLabels("aabbcc") == want);
}

int main() {
    testClassic();
    testAllUnique();
    testAllSame();
    testSingleChar();
    testMultipleEqualPartitions();
    std::printf("all tests passed\n");
    return 0;
}
