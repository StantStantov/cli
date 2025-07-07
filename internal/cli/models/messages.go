package models

type AuthSuccessMsg struct {
	Token    string
	Username string
}

type LogoutMsg struct{}

type OpenChatMsg struct{}

type UsernameChangeMsg struct {
	NewUsername string
}

type ChatKeyHandledMsg struct{}

type ChatClosedMsg struct{}
