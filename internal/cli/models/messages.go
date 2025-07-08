package models

type LoginWithOAuthMsg struct {
	Provider string
}

type RegisterWithPasswordMsg struct {
	Login    string
	Password string
}

type OAuthStartMsg struct {
	Provider string
	IsLogin  bool
}

type OAuthURIMsg struct {
	URI      string
	Provider string
	IsLogin  bool
}

type OAuthStatusCheckMsg struct {
	Provider string
}

type OAuthCancelMsg struct{}

type OAuthTimeoutMsg struct{}

type OAuthPollingResultMsg struct {
	Username string
	Gold     int
	Error    string
}

type AuthSuccessMsg struct {
	Username string
	Gold     int
}

type LogoutMsg struct{}

type OpenChatMsg struct{}

type UsernameChangeMsg struct {
	NewUsername string
	Gold        int
}

type ChatKeyHandledMsg struct{}

type ChatClosedMsg struct{}
