package models

import (
	"fmt"
	"lesta-start-battleship/cli/internal/cli/ui"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const matchmakingUrl = "ws://37.9.53.32:80/matchmaking/%s"

func formatMatchmakingUrl(matchType string) string {
	return fmt.Sprintf(matchmakingUrl, matchType)
}

type MatchmakingModel struct {
	username string

	selected int
}

func NewMatchmakingModel(username string) *MatchmakingModel {
	return &MatchmakingModel{
		username: username,
	}
}

func (m *MatchmakingModel) Init() tea.Cmd {
	return nil
}

const matchTypesAmount = 4

func (m *MatchmakingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				model := NewMatchmakingWaitScreenModel(m.username, "random")
				return model, model.Init()
			case 1:
				model := NewMatchmakingWaitScreenModel(m.username, "ranked")
				return model, model.Init()
			case 2:
				return m, nil
			case 3:
				model := NewMatchmakingCustomMenuModel(m.username)
				return model, model.Init()
			}
			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(0, m.username, 0, nil), nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *MatchmakingModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	menuItems := []string{
		"Случайный",
		"Рейтинговый",
		"Гильдейский",
		"Кастомный",
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
