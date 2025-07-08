package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

type GuildModel struct {
	id       int
	username string
	gold     int
	Member   *guilds.MemberResponse
	Guild    *guilds.GuildResponse
	selected int
	loading  bool
	errorMsg string
	Clients  *clientdeps.Client
}

func NewGuildModel(id int, username string, gold int, member *guilds.MemberResponse, guild *guilds.GuildResponse, clients *clientdeps.Client) *GuildModel {
	return &GuildModel{
		id:       id,
		username: username,
		gold:     gold,
		Member:   member,
		Guild:    guild,
		Clients:  clients,
	}
}

func (m *GuildModel) Init() tea.Cmd {
	return nil
}

func (m *GuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			return NewMainMenuModel(m.id, m.username, m.gold, m.Clients), nil
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

	if m.Member == nil || m.Guild == nil {
		sb.WriteString("\nВы не состоите в гильдии\n")
	} else {
		sb.WriteString(fmt.Sprintf("\nГильдия: [%s] %s\n", m.Guild.Tag, m.Guild.Title))
		sb.WriteString(fmt.Sprintf("Ваша роль: %s\n", m.Member.Role.Title))
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

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg + "\n"))
	}

	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - назад"))

	return sb.String()
}

func (m *GuildModel) getMenuItems() []string {
	if m.Member == nil || m.Guild == nil {
		return []string{"Вступить в гильдию", "Создать гильдию", "Список гильдий"}
	}

	switch m.Member.Role.Title {
	case "cabin_boy":
		return []string{"Список гильдий", "Список участников", "Чат гильдии", "Покинуть гильдию"}
	case "owner":
		return []string{"Объявить войну", "Запросы на войну", "Изменить гильдию", "Список участников", "Список гильдий",
			"Чат гильдии", "Запросы на вступление", "Удалить гильдию"}
	default:
		return []string{"Список гильдий", "Список участников", "Чат гильдии", "Запросы на вступление", "Покинуть гильдию"}
	}
}

func (m *GuildModel) handleMenuSelection() (tea.Model, tea.Cmd) {
	menuItems := m.getMenuItems()
	if len(menuItems) == 0 {
		return m, nil
	}

	selectedItem := menuItems[m.selected]

	switch selectedItem {
	case "Список гильдий":
		return NewGuildListModel(m, m.id, m.username, m.Clients), nil
	case "Список участников":
		return NewMembersListModel(m, m.id, m.username, m.Member.Role.Title, m.Guild.Tag, m.Guild.Title, m.Clients), nil
	case "Чат гильдии":
		// Инициализация чата гильдии с правильным guildID
		guildID := 0
		if m.Guild != nil {
			guildID = m.Guild.ID
		}
		return m, func() tea.Msg {
			return OpenChatMsg{
				GuildID: guildID,
			}
		}
	case "Покинуть гильдию":
		return NewExitGuildModel(m, m.id, m.username, m.gold, m.Guild.Tag, m.Guild.Title, m.Clients), nil
	case "Объявить войну":
		return m, nil
	case "Запросы на войну":
		return m, nil
	case "Изменить гильдию":
		return m, nil
	case "Запросы на вступление":
		return m, nil
	case "Удалить гильдию":
		return m, nil
	case "Создать гильдию":
		return m, nil
	case "Вступить в гильдию":
		return m, nil
	default:
		return m, nil
	}
}
