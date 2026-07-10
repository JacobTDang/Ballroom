#include <cassert>
#include <cstdio>
#include <vector>

#include "solution.hpp"

int main() {
    {
        Twitter tw;
        tw.postTweet(1, 5);
        assert((tw.getNewsFeed(1) == std::vector<int>{5}));

        tw.follow(1, 2);
        tw.postTweet(2, 6);
        assert((tw.getNewsFeed(1) == std::vector<int>{6, 5}));

        tw.unfollow(1, 2);
        assert((tw.getNewsFeed(1) == std::vector<int>{5}));
    }
    {
        Twitter tw;
        for (int i = 0; i < 15; i++) tw.postTweet(1, i);
        auto feed = tw.getNewsFeed(1);
        assert(feed.size() == 10);
        assert((feed == std::vector<int>{14, 13, 12, 11, 10, 9, 8, 7, 6, 5}));
    }
    {
        Twitter tw;
        tw.follow(1, 1);
        tw.postTweet(1, 100);
        assert((tw.getNewsFeed(1) == std::vector<int>{100}));
    }
    {
        Twitter tw;
        assert(tw.getNewsFeed(999).empty());
    }
    printf("all assertions passed\n");
    return 0;
}
