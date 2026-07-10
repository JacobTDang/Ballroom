class Twitter:
    """Simplified Twitter supporting posting, following, and a
    per-user news feed of the 10 most recent tweets from the user
    and everyone they follow."""

    def __init__(self):
        pass

    def post_tweet(self, user_id: int, tweet_id: int) -> None:
        raise NotImplementedError

    def get_news_feed(self, user_id: int) -> list[int]:
        raise NotImplementedError

    def follow(self, follower_id: int, followee_id: int) -> None:
        raise NotImplementedError

    def unfollow(self, follower_id: int, followee_id: int) -> None:
        raise NotImplementedError
