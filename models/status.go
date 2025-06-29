package models

type FriendshipStatus string

const (
	FriendshipRequested FriendshipStatus = "requested"
	FriendshipAccepted  FriendshipStatus = "accepted"
	FriendshipRejected  FriendshipStatus = "rejected"
)
