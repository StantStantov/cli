package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/cli/handlers"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

type GuildModel struct {
	username  string
	gold      int
	GuildInfo handlers.GuildResponse
	selected  int
	loading   bool
	err       error
	Clients   *clientdeps.Client
}

func NewGuildModel(username string, gold int, guildInfo handlers.GuildResponse, clients *clientdeps.Client) *GuildModel {
	return &GuildModel{
		username:  username,
		gold:      gold,
		GuildInfo: guildInfo,
		Clients:   clients,
		//loading:  true,
	}
}

func (m *GuildModel) Init() tea.Cmd {
	return nil
}

func (m *GuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case handlers.GuildResponse:
		m.GuildInfo = msg
		m.loading = false
		return m, nil

	case error:
		m.err = msg
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.Type {
		case tea.KeyUp:
			menuItems := m.getMenuItems()
			m.selected = (m.selected - 1 + len(menuItems)) % len(menuItems)
			return m, nil

		case tea.KeyDown:
			menuItems := m.getMenuItems()
			m.selected = (m.selected + 1) % len(menuItems)
			return m, nil

		case tea.KeyEnter:
			return m.handleMenuSelection()

		case tea.KeyEsc:
			return NewMainMenuModel(m.username, m.gold, m.Clients), nil
		}
	}

	return m, nil
}

func (m *GuildModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Гильдия"))
	sb.WriteString("\n")

	if m.loading {
		sb.WriteString("\nЗагрузка данных гильдии...")
		return sb.String()
	}

	if m.err != nil {
		sb.WriteString(ui.ErrorStyle.Render("\nОшибка: " + m.err.Error()))
		return sb.String()
	}

	if !m.GuildInfo.Member {
		sb.WriteString("\nВы не состоите в гильдии\n")
	} else {
		sb.WriteString(fmt.Sprintf("\nГильдия: [%s] %s\n", m.GuildInfo.Info.Tag, m.GuildInfo.Info.Name))
	}

	sb.WriteString("\n")
	menuItems := m.getMenuItems()
	for i, item := range menuItems {
		if i == m.selected {
			sb.WriteString(ui.SelectedStyle.Render("> " + item))
		} else {
			sb.WriteString(" " + item)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - назад"))

	return sb.String()
}

func (m *GuildModel) getMenuItems() []string {
	if !m.GuildInfo.Member {
		return []string{"Создать гильдию", "Вступить в гильдию", "Список гильдий"}
	}

	if m.GuildInfo.Owner {
		return []string{"Объявить войну", "Изменить роли", "Список участников", "Чат гильдии", "Удалить гильдию"}
	}

	return []string{"Список участников", "Чат гильдии", "Покинуть гильдию"}
}

func (m *GuildModel) handleMenuSelection() (tea.Model, tea.Cmd) {
	menuItems := m.getMenuItems()
	if len(menuItems) == 0 {
		return m, nil
	}

	selectedItem := menuItems[m.selected]

	switch selectedItem {
	case "Чат гильдии":
		//нужна реализация
		if m.GuildInfo.Member {
			return m, func() tea.Msg { return OpenChatMsg{} }
		}
		return m, nil
	case "Создать гильдию":
		return NewCreateGuildModel(m.username, m.gold, m.Clients), nil
	case "Вступить в гильдию":
		return m, nil
		//return NewJoinGuildModel(m.Username), nil
	// ... другие case
	default:
		return m, nil

	}
}

func (m *GuildModel) loadGuildData() tea.Msg {
	response, err := handlers.GetGuildInfo("dummy_token_" + m.username)
	if err != nil {
		return err
	}
	return response
}
