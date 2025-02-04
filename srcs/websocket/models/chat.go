package models

type Event struct {
	Type string `json:"type"`
}

type MessageEvent struct {
	Event
	Data       string `json:"data"`
	SenderID   uint64 `json:"senderId"`
	ReceiverID uint64 `json:"receiverId"`
}

type OnlineUsersEvent struct {
	Event
	Users []uint64 `json:"usersOnline"`
}

type UserStatusEvent struct {
	Event
	User uint64 `json:"user"`
}
