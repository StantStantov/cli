package models

import (
	"lesta-start-battleship/cli/internal/api/guilds"
)

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
	ID       int
	Username string
	Gold     int
	Error    string
}

type AuthSuccessMsg struct {
	ID       int
	Username string
	Gold     int
}

type LogoutMsg struct{}

type OpenChatMsg struct {
	GuildID int
}

type UsernameChangeMsg struct {
	NewUsername string
	Gold        int
}

type ChatKeyHandledMsg struct{}

type ChatClosedMsg struct{}

type GuildDataMsg struct {
	Member *guilds.MemberResponse
	Guild  *guilds.GuildResponse
}

type GuildNoMemberMsg struct{}

type MemberRoleChangeMsg struct {
	Username string
}

type MemberDeleteMsg struct {
	Username string
}

type GuildExitedMsg struct{}

type RequestProcessedMsg struct {
	Message string
}
