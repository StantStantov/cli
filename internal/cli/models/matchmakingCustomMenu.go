package models

import (
	"lesta-battleship/cli/internal/cli/ui"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MatchmakingCustomMenuModel struct {
	username string

	selected int
}

func NewMatchmakingCustomMenuModel(username string) *MatchmakingCustomMenuModel {
	return &MatchmakingCustomMenuModel{
		username: username,
	}
}

func (m *MatchmakingCustomMenuModel) Init() tea.Cmd {
	return nil
}

func (m *MatchmakingCustomMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.selected = (m.selected - 1 + matchTypesAmount) % matchTypesAmount
			return m, nil

		case tea.KeyDown:
			m.selected = (m.selected + 1) % matchTypesAmount
			return m, nil

		case tea.KeyEnter:
			switch m.selected {
			case 0:
				model := NewMatchmakingCustomRoomModel(m.username)
				return model, model.Init()
			case 1:
				model := NewMatchmakingCustomJoinModel(m.username)
				return model, model.Init()
			}
			return m, nil

		case tea.KeyEsc:
			return NewMatchmakingModel(m.username), nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *MatchmakingCustomMenuModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	menuItems := []string{
		"Создать",
		"Присоединиться",
	}

	for i, item := range menuItems {
		if i == m.selected {
			sb.WriteString(ui.SelectedStyle.Render("> " + item))
		} else {
			sb.WriteString(ui.NormalStyle.Render("  " + item))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - выход"))

	return sb.String()
}
