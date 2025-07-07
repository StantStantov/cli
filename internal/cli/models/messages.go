package models

type LoginWithPasswordMsg struct {
	Login    string
	Password string
}

type LoginWithOAuthMsg struct {
	Provider string
}

type RegisterWithPasswordMsg struct {
	Login    string
	Password string
}

type RegisterWithOAuthMsg struct {
	Provider string
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

type AuthSuccessMsg struct {
	Username string
	Gold     int
}

type LogoutMsg struct{}

type OpenChatMsg struct{}

type UsernameChangeMsg struct {
	NewUsername string
}

type ChatKeyHandledMsg struct{}

type ChatClosedMsg struct{}
