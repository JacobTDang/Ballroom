package main

import (
	"reflect"
	"testing"
)

func befriend(s *ChatServer, a, b string) {
	s.AddUser(a)
	s.AddUser(b)
	s.SendFriendRequest(a, b)
	s.AcceptFriendRequest(a, b)
}

func TestAddUserSucceedsOnce(t *testing.T) {
	s := NewChatServer()
	if !s.AddUser("alice") {
		t.Error("first AddUser = false, want true")
	}
	if s.AddUser("alice") {
		t.Error("second AddUser = true, want false")
	}
}

func TestRequestAndAcceptMakesFriendsBothWays(t *testing.T) {
	s := NewChatServer()
	befriend(s, "alice", "bob")
	if !s.AreFriends("alice", "bob") || !s.AreFriends("bob", "alice") {
		t.Error("friendship must be symmetric after accept")
	}
}

func TestRequestToUnknownUserFails(t *testing.T) {
	s := NewChatServer()
	s.AddUser("alice")
	if s.SendFriendRequest("alice", "ghost") {
		t.Error("request to unknown user = true, want false")
	}
	if s.SendFriendRequest("ghost", "alice") {
		t.Error("request from unknown user = true, want false")
	}
}

func TestSelfRequestFails(t *testing.T) {
	s := NewChatServer()
	s.AddUser("alice")
	if s.SendFriendRequest("alice", "alice") {
		t.Error("self-request = true, want false")
	}
}

func TestDuplicatePendingRequestFails(t *testing.T) {
	s := NewChatServer()
	s.AddUser("alice")
	s.AddUser("bob")
	if !s.SendFriendRequest("alice", "bob") {
		t.Fatal("first request = false, want true")
	}
	if s.SendFriendRequest("alice", "bob") {
		t.Error("duplicate request = true, want false")
	}
}

func TestRequestBetweenExistingFriendsFails(t *testing.T) {
	s := NewChatServer()
	befriend(s, "alice", "bob")
	if s.SendFriendRequest("alice", "bob") || s.SendFriendRequest("bob", "alice") {
		t.Error("request between friends = true, want false")
	}
}

func TestAcceptWithoutPendingRequestFails(t *testing.T) {
	s := NewChatServer()
	s.AddUser("alice")
	s.AddUser("bob")
	if s.AcceptFriendRequest("alice", "bob") {
		t.Error("accept with no pending request = true, want false")
	}
}

func TestAcceptIsDirectional(t *testing.T) {
	s := NewChatServer()
	s.AddUser("alice")
	s.AddUser("bob")
	s.SendFriendRequest("alice", "bob")
	if s.AcceptFriendRequest("bob", "alice") {
		t.Error("accepting the reversed direction = true, want false")
	}
	if !s.AcceptFriendRequest("alice", "bob") {
		t.Error("accepting the real request = false, want true")
	}
}

func TestAcceptConsumesTheRequest(t *testing.T) {
	s := NewChatServer()
	befriend(s, "alice", "bob")
	if s.AcceptFriendRequest("alice", "bob") {
		t.Error("second accept = true, want false -- the request was consumed")
	}
}

func TestMessageRequiresFriendship(t *testing.T) {
	s := NewChatServer()
	s.AddUser("alice")
	s.AddUser("bob")
	if s.SendMessage("alice", "bob", "hi") {
		t.Error("message between non-friends = true, want false")
	}
	s.SendFriendRequest("alice", "bob")
	s.AcceptFriendRequest("alice", "bob")
	if !s.SendMessage("alice", "bob", "hi") {
		t.Error("message between friends = false, want true")
	}
}

func TestReadMessagesReturnsFormattedAndDrains(t *testing.T) {
	s := NewChatServer()
	befriend(s, "alice", "bob")
	s.SendMessage("alice", "bob", "hello")
	s.SendMessage("alice", "bob", "there")
	got := s.ReadMessages("bob")
	want := []string{"alice: hello", "alice: there"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadMessages = %v, want %v", got, want)
	}
	if got := s.ReadMessages("bob"); len(got) != 0 {
		t.Errorf("second ReadMessages = %v, want empty -- reading drains", got)
	}
}

func TestMessagesFromMultipleSendersKeepArrivalOrder(t *testing.T) {
	s := NewChatServer()
	befriend(s, "alice", "carol")
	befriend(s, "bob", "carol")
	s.SendMessage("alice", "carol", "one")
	s.SendMessage("bob", "carol", "two")
	s.SendMessage("alice", "carol", "three")
	got := s.ReadMessages("carol")
	want := []string{"alice: one", "bob: two", "alice: three"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadMessages = %v, want %v", got, want)
	}
}

func TestReadMessagesForUnknownUserIsEmpty(t *testing.T) {
	if got := NewChatServer().ReadMessages("ghost"); len(got) != 0 {
		t.Errorf("ReadMessages(ghost) = %v, want empty", got)
	}
}
