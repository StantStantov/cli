package models

import (
	"lesta-battleship/cli/internal/cli/ui"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MatchmakingCustomJoinModel struct {
	username string
}

func NewMatchmakingCustomJoinModel(username string) *MatchmakingCustomJoinModel {
	return &MatchmakingCustomJoinModel{
		username: username,
	}
}

func (m *MatchmakingCustomJoinModel) Init() tea.Cmd {
	return tickEvery()
}

func (m *MatchmakingCustomJoinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *MatchmakingCustomJoinModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	sb.WriteString("Введите ID:")

	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - выход"))

	return sb.String()
}
