package main

// ChatServer is the in-memory chat core: users, friend requests,
// friendships, and inboxes. Messages only flow between friends.
type ChatServer struct {
	users   map[string]bool
	pending map[[2]string]bool // [from, to]
	friends map[[2]string]bool // sorted pair
	inboxes map[string][]string
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		users:   map[string]bool{},
		pending: map[[2]string]bool{},
		friends: map[[2]string]bool{},
		inboxes: map[string][]string{},
	}
}

func pair(a, b string) [2]string {
	if a > b {
		a, b = b, a
	}
	return [2]string{a, b}
}

// AddUser registers a user; false if the name is already taken.
func (s *ChatServer) AddUser(name string) bool {
	if s.users[name] {
		return false
	}
	s.users[name] = true
	return true
}

// SendFriendRequest sends a friend request. False if either user is
// unknown, the users are the same or already friends, or an identical
// request is already pending.
func (s *ChatServer) SendFriendRequest(fromUser, toUser string) bool {
	if !s.users[fromUser] || !s.users[toUser] || fromUser == toUser {
		return false
	}
	if s.friends[pair(fromUser, toUser)] {
		return false
	}
	req := [2]string{fromUser, toUser}
	if s.pending[req] {
		return false
	}
	s.pending[req] = true
	return true
}

// AcceptFriendRequest accepts the pending request fromUser sent to
// toUser. False if no such request is pending.
func (s *ChatServer) AcceptFriendRequest(fromUser, toUser string) bool {
	req := [2]string{fromUser, toUser}
	if !s.pending[req] {
		return false
	}
	delete(s.pending, req)
	s.friends[pair(fromUser, toUser)] = true
	return true
}

// AreFriends reports whether a and b are friends (symmetric).
func (s *ChatServer) AreFriends(a, b string) bool {
	return s.friends[pair(a, b)]
}

// SendMessage delivers text to toUser's inbox. False unless both users
// exist and are friends.
func (s *ChatServer) SendMessage(fromUser, toUser, text string) bool {
	if !s.AreFriends(fromUser, toUser) {
		return false
	}
	s.inboxes[toUser] = append(s.inboxes[toUser], fromUser+": "+text)
	return true
}

// ReadMessages returns and clears user's inbox, oldest first, each
// entry formatted "from: text".
func (s *ChatServer) ReadMessages(user string) []string {
	inbox := s.inboxes[user]
	s.inboxes[user] = nil
	return inbox
}
