from solution import ChatServer


def befriend(s, a, b):
    s.add_user(a)
    s.add_user(b)
    s.send_friend_request(a, b)
    s.accept_friend_request(a, b)


def test_add_user_succeeds_once():
    s = ChatServer()
    assert s.add_user("alice") is True
    assert s.add_user("alice") is False


def test_request_and_accept_makes_friends_both_ways():
    s = ChatServer()
    befriend(s, "alice", "bob")
    assert s.are_friends("alice", "bob") is True
    assert s.are_friends("bob", "alice") is True


def test_request_to_unknown_user_fails():
    s = ChatServer()
    s.add_user("alice")
    assert s.send_friend_request("alice", "ghost") is False
    assert s.send_friend_request("ghost", "alice") is False


def test_self_request_fails():
    s = ChatServer()
    s.add_user("alice")
    assert s.send_friend_request("alice", "alice") is False


def test_duplicate_pending_request_fails():
    s = ChatServer()
    s.add_user("alice")
    s.add_user("bob")
    assert s.send_friend_request("alice", "bob") is True
    assert s.send_friend_request("alice", "bob") is False


def test_request_between_existing_friends_fails():
    s = ChatServer()
    befriend(s, "alice", "bob")
    assert s.send_friend_request("alice", "bob") is False
    assert s.send_friend_request("bob", "alice") is False


def test_accept_without_pending_request_fails():
    s = ChatServer()
    s.add_user("alice")
    s.add_user("bob")
    assert s.accept_friend_request("alice", "bob") is False


def test_accept_is_directional():
    s = ChatServer()
    s.add_user("alice")
    s.add_user("bob")
    s.send_friend_request("alice", "bob")
    assert s.accept_friend_request("bob", "alice") is False
    assert s.accept_friend_request("alice", "bob") is True


def test_accept_consumes_the_request():
    s = ChatServer()
    befriend(s, "alice", "bob")
    assert s.accept_friend_request("alice", "bob") is False


def test_message_requires_friendship():
    s = ChatServer()
    s.add_user("alice")
    s.add_user("bob")
    assert s.send_message("alice", "bob", "hi") is False
    befriend2 = s.send_friend_request("alice", "bob") and s.accept_friend_request("alice", "bob")
    assert befriend2 is True
    assert s.send_message("alice", "bob", "hi") is True


def test_read_messages_returns_formatted_and_drains():
    s = ChatServer()
    befriend(s, "alice", "bob")
    s.send_message("alice", "bob", "hello")
    s.send_message("alice", "bob", "there")
    assert s.read_messages("bob") == ["alice: hello", "alice: there"]
    assert s.read_messages("bob") == []


def test_messages_from_multiple_senders_keep_arrival_order():
    s = ChatServer()
    befriend(s, "alice", "carol")
    befriend(s, "bob", "carol")
    s.send_message("alice", "carol", "one")
    s.send_message("bob", "carol", "two")
    s.send_message("alice", "carol", "three")
    assert s.read_messages("carol") == ["alice: one", "bob: two", "alice: three"]


def test_read_messages_for_unknown_user_is_empty():
    s = ChatServer()
    assert s.read_messages("ghost") == []
