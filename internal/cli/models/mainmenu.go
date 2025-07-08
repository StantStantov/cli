package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"lesta-battleship/cli/internal/cli/handlers"
	"lesta-battleship/cli/internal/cli/ui"
	"strings"
)

type MainMenuModel struct {
	username string
	selected int
}

func NewMainMenuModel(username string) *MainMenuModel {
	return &MainMenuModel{
		username: username,
		selected: 0,
	}
}

func (m *MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m *MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.selected = (m.selected - 1 + 6) % 6
			return m, nil

		case tea.KeyDown:
			m.selected = (m.selected + 1) % 6
			return m, nil

		case tea.KeyEnter:
			switch m.selected {
			case 0: // Бой
				return m, nil
			case 1: // Инвентарь
				return m, m.loadHandler
			case 2: // Магазин
				return m, m.loadHandler
			case 3: // Гильдия
				return m, m.loadHandler
			case 4: // Редактирование профиля
				return NewEditProfileModel(m.username), nil
			case 5: // Рейтинг
				return m, m.loadHandler
			}
			return m, nil

		case tea.KeyEsc:
			return m, func() tea.Msg { return LogoutMsg{} }

		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case handlers.InventoryResponse:
		return NewInventoryModel(m.username, msg), nil

	case handlers.ShopResponse:
		return NewShopModel(m.username, msg), nil

	case handlers.GuildResponse:
		return NewGuildModel(m.username, msg), nil

	case handlers.PlayerStats:
		return NewScoreboardModel(m.username, msg), nil
	}

	return m, nil
}

func (m *MainMenuModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	menuItems := []string{
		"⚔️  Бой",
		"🎒 Инвентарь",
		"🏪 Магазин",
		"🏰 Гильдия",
		"👤 Редактирование профиля",
		"🏆 Рейтинги",
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

func (m *MainMenuModel) loadHandler() tea.Msg {
	token := "dummy_token_" + m.username
	switch m.selected {
	case 1:
		items, err := handlers.InventoryHandler(token)
		if err != nil {
			return err
		}
		return items
	case 2:
		items, err := handlers.ItemsHandler(token)
		if err != nil {
			return err
		}
		return items
	case 3:
		guildInfo, err := handlers.GetGuildInfo(token)
		if err != nil {
			return err
		}
		return guildInfo
	case 5:
		stats, err := handlers.MyStatsHandler(token)
		if err != nil {
			return err
		}
		return stats
	}
	return nil
}
