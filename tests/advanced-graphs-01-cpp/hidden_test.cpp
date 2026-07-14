#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<std::string> FindItinerary(std::vector<std::vector<std::string>>& tickets);

void testDeadEnd() {
    std::vector<std::vector<std::string>> tickets = {
        {"JFK", "SFO"}, {"JFK", "ATL"}, {"SFO", "ATL"}, {"ATL", "JFK"}, {"ATL", "SFO"}
    };
    std::vector<std::string> want = {"JFK", "ATL", "JFK", "SFO", "ATL", "SFO"};
    assert(FindItinerary(tickets) == want);
}

void testSimple() {
    std::vector<std::vector<std::string>> tickets = {
        {"MUC", "LHR"}, {"JFK", "MUC"}, {"SFO", "SJC"}, {"LHR", "SFO"}
    };
    std::vector<std::string> want = {"JFK", "MUC", "LHR", "SFO", "SJC"};
    assert(FindItinerary(tickets) == want);
}

void testLexicalTieBreak() {
    std::vector<std::vector<std::string>> tickets = {
        {"JFK", "KUL"}, {"JFK", "NRT"}, {"NRT", "JFK"}
    };
    std::vector<std::string> want = {"JFK", "NRT", "JFK", "KUL"};
    assert(FindItinerary(tickets) == want);
}

void testSimpleTwoCycle() {
    std::vector<std::vector<std::string>> tickets = {{"JFK", "A"}, {"A", "JFK"}};
    std::vector<std::string> want = {"JFK", "A", "JFK"};
    assert(FindItinerary(tickets) == want);
}

void testBranchingAtOrigin() {
    std::vector<std::vector<std::string>> tickets = {
        {"JFK", "B"}, {"JFK", "A"}, {"B", "JFK"}, {"A", "JFK"}
    };
    std::vector<std::string> want = {"JFK", "A", "JFK", "B", "JFK"};
    assert(FindItinerary(tickets) == want);
}

int main() {
    testDeadEnd();
    testSimple();
    testLexicalTieBreak();
    testSimpleTwoCycle();
    testBranchingAtOrigin();
    std::printf("all tests passed\n");
    return 0;
}
