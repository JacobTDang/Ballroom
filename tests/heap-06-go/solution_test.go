package main

import (
	"reflect"
	"testing"
)

func TestTwitter(t *testing.T) {
	tw := NewTwitter()
	tw.PostTweet(1, 5)
	if got := tw.GetNewsFeed(1); !reflect.DeepEqual(got, []int{5}) {
		t.Fatalf("GetNewsFeed(1) = %v, want [5]", got)
	}

	tw.Follow(1, 2)
	tw.PostTweet(2, 6)
	if got := tw.GetNewsFeed(1); !reflect.DeepEqual(got, []int{6, 5}) {
		t.Fatalf("GetNewsFeed(1) = %v, want [6 5]", got)
	}

	tw.Unfollow(1, 2)
	if got := tw.GetNewsFeed(1); !reflect.DeepEqual(got, []int{5}) {
		t.Fatalf("GetNewsFeed(1) = %v, want [5]", got)
	}
}

func TestTwitter_NewsFeedCapsAtTenMostRecent(t *testing.T) {
	tw := NewTwitter()
	for i := 0; i < 15; i++ {
		tw.PostTweet(1, i)
	}
	got := tw.GetNewsFeed(1)
	if len(got) != 10 {
		t.Fatalf("GetNewsFeed(1) length = %d, want 10", len(got))
	}
	want := []int{14, 13, 12, 11, 10, 9, 8, 7, 6, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetNewsFeed(1) = %v, want %v", got, want)
	}
}

func TestTwitter_SelfFollowIsNoop(t *testing.T) {
	tw := NewTwitter()
	tw.Follow(1, 1)
	tw.PostTweet(1, 100)
	got := tw.GetNewsFeed(1)
	if !reflect.DeepEqual(got, []int{100}) {
		t.Errorf("GetNewsFeed(1) = %v, want [100] (self-follow shouldn't duplicate)", got)
	}
}

func TestTwitter_EmptyFeedForUnknownUser(t *testing.T) {
	tw := NewTwitter()
	got := tw.GetNewsFeed(999)
	if len(got) != 0 {
		t.Errorf("GetNewsFeed(999) = %v, want empty", got)
	}
}
