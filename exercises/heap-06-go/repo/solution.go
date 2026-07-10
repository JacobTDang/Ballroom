package main

// Twitter is a simplified Twitter supporting posting, following, and
// a per-user news feed of the 10 most recent tweets from the user
// and everyone they follow.
type Twitter struct {
}

func NewTwitter() *Twitter {
	return &Twitter{}
}

func (tw *Twitter) PostTweet(userId, tweetId int) {
	// TODO: implement
}

func (tw *Twitter) GetNewsFeed(userId int) []int {
	// TODO: implement
	return nil
}

func (tw *Twitter) Follow(followerId, followeeId int) {
	// TODO: implement
}

func (tw *Twitter) Unfollow(followerId, followeeId int) {
	// TODO: implement
}
