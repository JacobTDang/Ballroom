package main

import "sort"

// tweet pairs a tweet id with the global post order it was made in,
// so feeds across multiple users can be merged and sorted by recency.
type tweet struct {
	id    int
	order int
}

// Twitter is a simplified Twitter supporting posting, following, and
// a per-user news feed of the 10 most recent tweets from the user
// and everyone they follow.
type Twitter struct {
	order     int
	tweets    map[int][]tweet
	following map[int]map[int]bool
}

func NewTwitter() *Twitter {
	return &Twitter{
		tweets:    make(map[int][]tweet),
		following: make(map[int]map[int]bool),
	}
}

func (tw *Twitter) PostTweet(userId, tweetId int) {
	tw.tweets[userId] = append(tw.tweets[userId], tweet{id: tweetId, order: tw.order})
	tw.order++
}

func (tw *Twitter) GetNewsFeed(userId int) []int {
	users := map[int]bool{userId: true}
	for f := range tw.following[userId] {
		users[f] = true
	}

	var all []tweet
	for u := range users {
		all = append(all, tw.tweets[u]...)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].order > all[j].order })

	limit := 10
	if len(all) < limit {
		limit = len(all)
	}
	res := make([]int, limit)
	for i := 0; i < limit; i++ {
		res[i] = all[i].id
	}
	return res
}

func (tw *Twitter) Follow(followerId, followeeId int) {
	if followerId == followeeId {
		return
	}
	if tw.following[followerId] == nil {
		tw.following[followerId] = make(map[int]bool)
	}
	tw.following[followerId][followeeId] = true
}

func (tw *Twitter) Unfollow(followerId, followeeId int) {
	delete(tw.following[followerId], followeeId)
}
