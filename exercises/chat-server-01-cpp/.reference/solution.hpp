#pragma once

#include <map>
#include <set>
#include <string>
#include <utility>
#include <vector>

// In-memory chat core: users, friend requests, friendships, and
// inboxes. Messages only flow between friends.
class ChatServer {
public:
    ChatServer() {}

    // Register a user; false if the name is already taken.
    bool add_user(const std::string& name) {
        return users_.insert(name).second;
    }

    // Send a friend request. False if either user is unknown, the users
    // are the same or already friends, or an identical request is
    // already pending.
    bool send_friend_request(const std::string& from, const std::string& to) {
        if (!users_.count(from) || !users_.count(to) || from == to) return false;
        if (friends_.count(pair(from, to))) return false;
        return pending_.insert({from, to}).second;
    }

    // Accept the pending request `from` sent to `to`. False if no such
    // request is pending.
    bool accept_friend_request(const std::string& from, const std::string& to) {
        if (pending_.erase({from, to}) == 0) return false;
        friends_.insert(pair(from, to));
        return true;
    }

    // Symmetric friendship check.
    bool are_friends(const std::string& a, const std::string& b) {
        return friends_.count(pair(a, b)) > 0;
    }

    // Deliver text to `to`'s inbox. False unless both users exist and
    // are friends.
    bool send_message(const std::string& from, const std::string& to,
                      const std::string& text) {
        if (!are_friends(from, to)) return false;
        inboxes_[to].push_back(from + ": " + text);
        return true;
    }

    // Return and clear the user's inbox, oldest first, each entry
    // formatted "from: text".
    std::vector<std::string> read_messages(const std::string& user) {
        std::vector<std::string> inbox = std::move(inboxes_[user]);
        inboxes_[user].clear();
        return inbox;
    }

private:
    static std::pair<std::string, std::string> pair(const std::string& a,
                                                    const std::string& b) {
        return a < b ? std::make_pair(a, b) : std::make_pair(b, a);
    }

    std::set<std::string> users_;
    std::set<std::pair<std::string, std::string>> pending_;
    std::set<std::pair<std::string, std::string>> friends_;
    std::map<std::string, std::vector<std::string>> inboxes_;
};
