#pragma once

#include <string>
#include <vector>

// In-memory chat core: users, friend requests, friendships, and
// inboxes. Messages only flow between friends.
class ChatServer {
public:
    ChatServer() {}

    // Register a user; false if the name is already taken.
    bool add_user(const std::string& name) {
        // TODO: implement
        return false;
    }

    // Send a friend request. False if either user is unknown, the users
    // are the same or already friends, or an identical request is
    // already pending.
    bool send_friend_request(const std::string& from, const std::string& to) {
        // TODO: implement
        return false;
    }

    // Accept the pending request `from` sent to `to`. False if no such
    // request is pending.
    bool accept_friend_request(const std::string& from, const std::string& to) {
        // TODO: implement
        return false;
    }

    // Symmetric friendship check.
    bool are_friends(const std::string& a, const std::string& b) {
        // TODO: implement
        return false;
    }

    // Deliver text to `to`'s inbox. False unless both users exist and
    // are friends.
    bool send_message(const std::string& from, const std::string& to,
                      const std::string& text) {
        // TODO: implement
        return false;
    }

    // Return and clear the user's inbox, oldest first, each entry
    // formatted "from: text".
    std::vector<std::string> read_messages(const std::string& user) {
        // TODO: implement
        return {};
    }
};
