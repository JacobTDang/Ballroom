# Design a Chat Server

Design the in-memory core of a chat service: users, friendships built
through requests, and messages that can only flow between friends.
(No networking — model the state machine.)

- `add_user(name)` — register a user; false if the name is taken
- `send_friend_request(from, to)` — false if either user is unknown,
  `from == to`, they are already friends, or an identical request is
  already pending
- `accept_friend_request(from, to)` — accept the pending request that
  `from` sent to `to`; false if no such request is pending
- `are_friends(a, b)` — symmetric friendship check
- `send_message(from, to, text)` — false unless both exist and are
  friends; otherwise delivered to `to`'s inbox
- `read_messages(user)` — return and clear the user's inbox, in
  arrival order, each entry formatted `"from: text"`

## Examples

```
add_user("alice")            -> true
add_user("bob")              -> true
send_message("alice", "bob", "hi")  -> false   (not friends)
send_friend_request("alice", "bob") -> true
accept_friend_request("alice", "bob") -> true
send_message("alice", "bob", "hi")  -> true
read_messages("bob")         -> ["alice: hi"]
read_messages("bob")         -> []
```

## Constraints

- Friendship is symmetric; requests are directional until accepted.
- Accepting consumes the pending request.
