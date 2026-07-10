from solution import Twitter


def test_twitter():
    tw = Twitter()
    tw.post_tweet(1, 5)
    assert tw.get_news_feed(1) == [5]

    tw.follow(1, 2)
    tw.post_tweet(2, 6)
    assert tw.get_news_feed(1) == [6, 5]

    tw.unfollow(1, 2)
    assert tw.get_news_feed(1) == [5]


def test_news_feed_caps_at_ten_most_recent():
    tw = Twitter()
    for i in range(15):
        tw.post_tweet(1, i)
    feed = tw.get_news_feed(1)
    assert len(feed) == 10
    assert feed == [14, 13, 12, 11, 10, 9, 8, 7, 6, 5]


def test_self_follow_is_noop():
    tw = Twitter()
    tw.follow(1, 1)
    tw.post_tweet(1, 100)
    assert tw.get_news_feed(1) == [100]


def test_empty_feed_for_unknown_user():
    tw = Twitter()
    assert tw.get_news_feed(999) == []
