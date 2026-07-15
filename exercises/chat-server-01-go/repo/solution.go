package main

// ChatServer is the in-memory chat core: users, friend requests,
// friendships, and inboxes. Messages only flow between friends.
type ChatServer struct {
}

func NewChatServer() *ChatServer {
	return &ChatServer{}
}

// AddUser registers a user; false if the name is already taken.
func (s *ChatServer) AddUser(name string) bool {
	// TODO: implement
	return false
}

// SendFriendRequest sends a friend request. False if either user is
// unknown, the users are the same or already friends, or an identical
// request is already pending.
func (s *ChatServer) SendFriendRequest(fromUser, toUser string) bool {
	// TODO: implement
	return false
}

// AcceptFriendRequest accepts the pending request fromUser sent to
// toUser. False if no such request is pending.
func (s *ChatServer) AcceptFriendRequest(fromUser, toUser string) bool {
	// TODO: implement
	return false
}

// AreFriends reports whether a and b are friends (symmetric).
func (s *ChatServer) AreFriends(a, b string) bool {
	// TODO: implement
	return false
}

// SendMessage delivers text to toUser's inbox. False unless both users
// exist and are friends.
func (s *ChatServer) SendMessage(fromUser, toUser, text string) bool {
	// TODO: implement
	return false
}

// ReadMessages returns and clears user's inbox, oldest first, each
// entry formatted "from: text".
func (s *ChatServer) ReadMessages(user string) []string {
	// TODO: implement
	return nil
}
