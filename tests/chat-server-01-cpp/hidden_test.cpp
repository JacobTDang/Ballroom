#include <cassert>
#include <cstdio>
#include <vector>

#include "solution.hpp"

static void befriend(ChatServer& s, const std::string& a, const std::string& b) {
    s.add_user(a);
    s.add_user(b);
    s.send_friend_request(a, b);
    s.accept_friend_request(a, b);
}

int main() {
    {
        ChatServer s;
        assert(s.add_user("alice") == true);           // add once
        assert(s.add_user("alice") == false);
    }
    {
        ChatServer s;                                  // symmetric friendship
        befriend(s, "alice", "bob");
        assert(s.are_friends("alice", "bob"));
        assert(s.are_friends("bob", "alice"));
    }
    {
        ChatServer s;                                  // unknown users
        s.add_user("alice");
        assert(s.send_friend_request("alice", "ghost") == false);
        assert(s.send_friend_request("ghost", "alice") == false);
    }
    {
        ChatServer s;                                  // self request
        s.add_user("alice");
        assert(s.send_friend_request("alice", "alice") == false);
    }
    {
        ChatServer s;                                  // duplicate pending
        s.add_user("alice");
        s.add_user("bob");
        assert(s.send_friend_request("alice", "bob") == true);
        assert(s.send_friend_request("alice", "bob") == false);
    }
    {
        ChatServer s;                                  // already friends
        befriend(s, "alice", "bob");
        assert(s.send_friend_request("alice", "bob") == false);
        assert(s.send_friend_request("bob", "alice") == false);
    }
    {
        ChatServer s;                                  // accept without request
        s.add_user("alice");
        s.add_user("bob");
        assert(s.accept_friend_request("alice", "bob") == false);
    }
    {
        ChatServer s;                                  // accept is directional
        s.add_user("alice");
        s.add_user("bob");
        s.send_friend_request("alice", "bob");
        assert(s.accept_friend_request("bob", "alice") == false);
        assert(s.accept_friend_request("alice", "bob") == true);
    }
    {
        ChatServer s;                                  // accept consumes
        befriend(s, "alice", "bob");
        assert(s.accept_friend_request("alice", "bob") == false);
    }
    {
        ChatServer s;                                  // friendship gates messages
        s.add_user("alice");
        s.add_user("bob");
        assert(s.send_message("alice", "bob", "hi") == false);
        s.send_friend_request("alice", "bob");
        s.accept_friend_request("alice", "bob");
        assert(s.send_message("alice", "bob", "hi") == true);
    }
    {
        ChatServer s;                                  // read formats and drains
        befriend(s, "alice", "bob");
        s.send_message("alice", "bob", "hello");
        s.send_message("alice", "bob", "there");
        std::vector<std::string> want = {"alice: hello", "alice: there"};
        assert(s.read_messages("bob") == want);
        assert(s.read_messages("bob").empty());
    }
    {
        ChatServer s;                                  // arrival order across senders
        befriend(s, "alice", "carol");
        befriend(s, "bob", "carol");
        s.send_message("alice", "carol", "one");
        s.send_message("bob", "carol", "two");
        s.send_message("alice", "carol", "three");
        std::vector<std::string> want = {"alice: one", "bob: two", "alice: three"};
        assert(s.read_messages("carol") == want);
    }
    {
        ChatServer s;
        assert(s.read_messages("ghost").empty());      // unknown user
    }
    printf("all assertions passed\n");
    return 0;
}
