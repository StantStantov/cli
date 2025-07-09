package models

import (
	"lesta-battleship/cli/internal/cli/ui"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MatchmakingCustomRoomModel struct {
	username string
}

func NewMatchmakingCustomRoomModel(username string) *MatchmakingCustomRoomModel {
	return &MatchmakingCustomRoomModel{
		username: username,
	}
}

func (m *MatchmakingCustomRoomModel) Init() tea.Cmd {
	return nil
}

func (m *MatchmakingCustomRoomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewMatchmakingCustomMenuModel(m.username), nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *MatchmakingCustomRoomModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	sb.WriteString("Комната")

	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - выход"))

	return sb.String()
}
