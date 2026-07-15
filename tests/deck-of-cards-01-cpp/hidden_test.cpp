#include <cassert>
#include <cstdio>
#include <set>
#include <utility>

#include "solution.hpp"

static std::set<std::pair<std::string, int>> as_set(const std::vector<Card>& cards) {
    std::set<std::pair<std::string, int>> s;
    for (const auto& c : cards) s.insert({c.suit, c.rank});
    return s;
}

int main() {
    assert(Deck().size() == 52);                       // full new deck
    {
        Deck d;
        auto got = d.deal(1);
        assert(got.size() == 1 && got[0] == (Card{"clubs", 1}));
    }
    {
        Deck d;                                        // canonical order
        d.deal(1);
        auto got = d.deal(3);
        assert(got.size() == 3);
        assert(got[0] == (Card{"clubs", 2}));
        assert(got[1] == (Card{"clubs", 3}));
        assert(got[2] == (Card{"clubs", 4}));
    }
    {
        Deck d;                                        // suit boundary
        d.deal(13);
        auto got = d.deal(1);
        assert(got.size() == 1 && got[0] == (Card{"diamonds", 1}));
    }
    {
        Deck d;
        d.deal(5);
        assert(d.size() == 47);                        // deal reduces size
    }
    {
        Deck d;                                        // 52 unique cards
        assert(as_set(d.deal(52)).size() == 52);
    }
    {
        Deck d;                                        // deal past empty
        d.deal(50);
        assert(d.deal(5).size() == 2);
        assert(d.size() == 0);
        assert(d.deal(1).empty());
    }
    {
        Deck d;                                        // reset restores
        d.deal(30);
        d.reset();
        assert(d.size() == 52);
        auto got = d.deal(1);
        assert(got.size() == 1 && got[0] == (Card{"clubs", 1}));
    }
    {
        Deck d;                                        // shuffle keeps cards
        auto dealt = d.deal(10);
        d.shuffle();
        auto rest = d.deal(52);
        assert(rest.size() == 42);
        auto all = as_set(rest);
        for (const auto& c : dealt) all.insert({c.suit, c.rank});
        assert(all.size() == 52);
    }
    {
        Deck d;
        d.deal(7);
        d.shuffle();
        assert(d.size() == 45);                        // shuffle keeps size
    }
    printf("all assertions passed\n");
    return 0;
}
