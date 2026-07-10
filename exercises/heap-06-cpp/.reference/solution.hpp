#pragma once

#include <algorithm>
#include <unordered_map>
#include <unordered_set>
#include <vector>

// Twitter is a simplified Twitter supporting posting, following, and
// a per-user news feed of the 10 most recent tweets from the user
// and everyone they follow.
class Twitter {
public:
    void postTweet(int userId, int tweetId) {
        tweets_[userId].push_back({order_, tweetId});
        order_++;
    }

    std::vector<int> getNewsFeed(int userId) {
        std::unordered_set<int> users = {userId};
        for (int f : following_[userId]) users.insert(f);

        std::vector<std::pair<int, int>> all;  // (order, tweetId)
        for (int u : users) {
            for (auto& t : tweets_[u]) all.push_back(t);
        }
        std::sort(all.begin(), all.end(),
                  [](const auto& a, const auto& b) { return a.first > b.first; });

        std::vector<int> res;
        for (size_t i = 0; i < all.size() && i < 10; i++) {
            res.push_back(all[i].second);
        }
        return res;
    }

    void follow(int followerId, int followeeId) {
        if (followerId == followeeId) return;
        following_[followerId].insert(followeeId);
    }

    void unfollow(int followerId, int followeeId) {
        following_[followerId].erase(followeeId);
    }

private:
    int order_ = 0;
    std::unordered_map<int, std::vector<std::pair<int, int>>> tweets_;  // userId -> (order, tweetId)
    std::unordered_map<int, std::unordered_set<int>> following_;
};
