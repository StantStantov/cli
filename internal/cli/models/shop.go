package models

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/cli/handlers"
	"lesta-start-battleship/cli/internal/cli/ui"
	"strings"
)

type ShopModel struct {
	username string
	items    handlers.ShopResponse
	selected int
	category int // 0-предметы, 1-акции, 2-сундуки
	err      error
}

func NewShopModel(username string, items handlers.ShopResponse) *ShopModel {
	return &ShopModel{
		username: username,
		items:    items,
		selected: 0,
		category: 0,
	}
}

func (m *ShopModel) Init() tea.Cmd {
	return m.loadItems
}

func (m *ShopModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case handlers.ShopResponse:
		m.items = msg
		return m, nil

	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			m.category = (m.category - 1 + 3) % 3
			return m, m.loadItems

		case tea.KeyRight:
			m.category = (m.category + 1) % 3
			return m, m.loadItems

		case tea.KeyUp:
			if len(m.items.Items) > 0 {
				m.selected = (m.selected - 1 + len(m.items.Items)) % len(m.items.Items)
			}
			return m, nil

		case tea.KeyDown:
			if len(m.items.Items) > 0 {
				m.selected = (m.selected + 1) % len(m.items.Items)
			}
			return m, nil

		case tea.KeyEnter:
			// Здесь будет логика покупки

			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(m.username), nil
		}
	}
	return m, nil
}

func (m *ShopModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Магазин"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("Пользователь: %s					Balance: %d 💰", m.username, m.items.Balance)))
	sb.WriteString("\n\n")

	// Отображение категорий
	categories := []string{"Предметы", "Акции", "Сундуки"}
	for i, cat := range categories {
		if i == m.category {
			sb.WriteString(ui.SelectedStyle.Render("[" + cat + "] "))
		} else {
			sb.WriteString(ui.NormalStyle.Render(cat + " "))
		}
	}
	sb.WriteString("\n\n")

	if m.err != nil {
		sb.WriteString(ui.ErrorStyle.Render("Ошибка: " + m.err.Error()))
		return sb.String()
	}

	if len(m.items.Items) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Товары отсутствуют"))
		return sb.String()
	}

	for i, item := range m.items.Items {
		if i == m.selected {
			sb.WriteString(ui.SelectedStyle.Render("> " + item.Name))
		} else {
			sb.WriteString(ui.NormalStyle.Render("  " + item.Name))
		}
		sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf(" - %d %s", item.Price, item.Currency)))
		sb.WriteString("\n")
		sb.WriteString(ui.NormalStyle.Render("   " + item.Description))
		sb.WriteString("\n\n")
	}

	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("←/→ - переключение категорий, ↑/↓ - выбор, Enter - купить, Esc - назад"))

	return sb.String()
}

func (m *ShopModel) loadItems() tea.Msg {
	token := "dummy_token_" + m.username
	var items handlers.ShopResponse
	var err error

	switch m.category {
	case 0:
		items, err = handlers.ItemsHandler(token)
	case 1:
		items, err = handlers.PromoHandler(token)
	case 2:
		items, err = handlers.ChestsHandler(token)
	}

	if err != nil {
		return err
	}
	return items
}
