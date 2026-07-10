#pragma once

#include <unordered_map>
#include <unordered_set>
#include <vector>

// Twitter is a simplified Twitter supporting posting, following, and
// a per-user news feed of the 10 most recent tweets from the user
// and everyone they follow.
class Twitter {
public:
    void postTweet(int userId, int tweetId) {
        // TODO: implement
    }

    std::vector<int> getNewsFeed(int userId) {
        // TODO: implement
        return {};
    }

    void follow(int followerId, int followeeId) {
        // TODO: implement
    }

    void unfollow(int followerId, int followeeId) {
        // TODO: implement
    }
};
