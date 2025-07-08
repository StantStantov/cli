package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/cli/handlers"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

type CreateGuildModel struct {
	username    string
	gold        int
	name        string
	tag         string
	activeField int
	errorMsg    string
	Clients     *clientdeps.Client
}

func NewCreateGuildModel(username string, gold int, clients *clientdeps.Client) *CreateGuildModel {
	return &CreateGuildModel{
		username: username,
		gold:     gold,
		Clients:  clients,
	}
}

func (m *CreateGuildModel) Init() tea.Cmd {
	return nil
}

func (m *CreateGuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.activeField == 1 {
				//Заглушка для создания гильдии
				return NewGuildModel(m.username, m.gold, handlers.GuildResponse{
					Member: true,
					Owner:  true,
					Info: handlers.GuildInfo{
						Id:   1,
						Name: m.name,
						Tag:  m.tag,
					},
				}, m.Clients), nil
			}
			m.activeField = (m.activeField + 1) % 2
			return m, nil

		case tea.KeyBackspace:
			if m.activeField == 0 && len(m.name) > 0 {
				m.name = m.name[:len(m.name)-1]
			} else if m.activeField == 1 && len(m.tag) > 0 {
				m.tag = m.tag[:len(m.tag)-1]
			}
			return m, nil

		case tea.KeyRunes:
			if m.activeField == 0 {
				m.name += string(msg.Runes)
			} else {
				m.tag += string(msg.Runes)
			}
			return m, nil

		case tea.KeyEsc:
			return NewGuildModel(m.username, m.gold, handlers.GuildResponse{}, m.Clients), nil
		}
	}
	return m, nil
}

func (m *CreateGuildModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Создание гильдии"))
	sb.WriteString("\n\n")

	sb.WriteString("Название гильдии:\n")
	if m.activeField == 0 {
		sb.WriteString(ui.SelectedStyle.Render("> " + m.name + "_"))
	} else {
		sb.WriteString(" " + m.name)
	}
	sb.WriteString("\n\n")

	sb.WriteString("Тег гильдии:\n")
	if m.activeField == 1 {
		sb.WriteString(ui.SelectedStyle.Render("> " + m.tag + "_"))
	} else {
		sb.WriteString(" " + m.tag)
	}

	if m.errorMsg != "" {
		sb.WriteString("\n\n")
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg))
	}

	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Enter - подтвердить, Esc - назад"))

	return sb.String()
}
