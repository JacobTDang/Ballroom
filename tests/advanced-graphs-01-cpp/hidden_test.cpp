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

int main() {
    testDeadEnd();
    testSimple();
    testLexicalTieBreak();
    std::printf("all tests passed\n");
    return 0;
}
