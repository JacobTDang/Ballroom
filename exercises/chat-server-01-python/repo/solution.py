class ChatServer:
    """In-memory chat core: users, friend requests, friendships, and
    inboxes. Messages only flow between friends."""

    def __init__(self):
        pass

    def add_user(self, name: str) -> bool:
        """Register a user; False if the name is already taken."""
        raise NotImplementedError

    def send_friend_request(self, from_user: str, to_user: str) -> bool:
        """Send a friend request. False if either user is unknown, the
        users are the same or already friends, or an identical request
        is already pending."""
        raise NotImplementedError

    def accept_friend_request(self, from_user: str, to_user: str) -> bool:
        """Accept the pending request from_user sent to to_user. False
        if no such request is pending."""
        raise NotImplementedError

    def are_friends(self, a: str, b: str) -> bool:
        """Symmetric friendship check."""
        raise NotImplementedError

    def send_message(self, from_user: str, to_user: str, text: str) -> bool:
        """Deliver text to to_user's inbox. False unless both users
        exist and are friends."""
        raise NotImplementedError

    def read_messages(self, user: str) -> list:
        """Return and clear user's inbox, oldest first, each entry
        formatted "from: text"."""
        raise NotImplementedError
