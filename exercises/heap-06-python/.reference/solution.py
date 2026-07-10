class Twitter:
    """Simplified Twitter supporting posting, following, and a
    per-user news feed of the 10 most recent tweets from the user
    and everyone they follow."""

    def __init__(self):
        self.order = 0
        self.tweets: dict[int, list[tuple[int, int]]] = {}  # userId -> [(order, tweetId)]
        self.following: dict[int, set[int]] = {}

    def post_tweet(self, user_id: int, tweet_id: int) -> None:
        self.tweets.setdefault(user_id, []).append((self.order, tweet_id))
        self.order += 1

    def get_news_feed(self, user_id: int) -> list[int]:
        users = {user_id} | self.following.get(user_id, set())
        all_tweets: list[tuple[int, int]] = []
        for u in users:
            all_tweets.extend(self.tweets.get(u, []))
        all_tweets.sort(key=lambda t: t[0], reverse=True)
        return [tweet_id for _, tweet_id in all_tweets[:10]]

    def follow(self, follower_id: int, followee_id: int) -> None:
        if follower_id == followee_id:
            return
        self.following.setdefault(follower_id, set()).add(followee_id)

    def unfollow(self, follower_id: int, followee_id: int) -> None:
        self.following.get(follower_id, set()).discard(followee_id)
