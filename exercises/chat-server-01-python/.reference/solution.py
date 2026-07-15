class ChatServer:
    """In-memory chat core: users, friend requests, friendships, and
    inboxes. Messages only flow between friends."""

    def __init__(self):
        self._users = set()
        self._pending = set()   # (from_user, to_user)
        self._friends = set()   # frozenset({a, b})
        self._inboxes = {}      # user -> list of "from: text"

    def add_user(self, name: str) -> bool:
        if name in self._users:
            return False
        self._users.add(name)
        self._inboxes[name] = []
        return True

    def send_friend_request(self, from_user: str, to_user: str) -> bool:
        if from_user not in self._users or to_user not in self._users:
            return False
        if from_user == to_user:
            return False
        if frozenset((from_user, to_user)) in self._friends:
            return False
        if (from_user, to_user) in self._pending:
            return False
        self._pending.add((from_user, to_user))
        return True

    def accept_friend_request(self, from_user: str, to_user: str) -> bool:
        if (from_user, to_user) not in self._pending:
            return False
        self._pending.remove((from_user, to_user))
        self._friends.add(frozenset((from_user, to_user)))
        return True

    def are_friends(self, a: str, b: str) -> bool:
        return frozenset((a, b)) in self._friends

    def send_message(self, from_user: str, to_user: str, text: str) -> bool:
        if not self.are_friends(from_user, to_user):
            return False
        self._inboxes[to_user].append(f"{from_user}: {text}")
        return True

    def read_messages(self, user: str) -> list:
        if user not in self._users:
            return []
        inbox = self._inboxes[user]
        self._inboxes[user] = []
        return inbox
