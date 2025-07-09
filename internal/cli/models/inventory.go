package models

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"lesta-battleship/cli/internal/cli/handlers"
	"lesta-battleship/cli/internal/cli/ui"
	"strings"
)

type InventoryModel struct {
	username    string
	items       handlers.InventoryResponse
	selected    int
	showDetails bool
}

func NewInventoryModel(username string, items handlers.InventoryResponse) *InventoryModel {
	return &InventoryModel{
		username:    username,
		items:       items,
		selected:    0,
		showDetails: false,
	}
}

func (m *InventoryModel) Init() tea.Cmd {
	return nil
}

func (m *InventoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if !m.showDetails && len(m.items) > 0 {
				m.selected = (m.selected - 1 + len(m.items)) % len(m.items)
			}
			return m, nil

		case tea.KeyDown:
			if !m.showDetails && len(m.items) > 0 {
				m.selected = (m.selected + 1) % len(m.items)
			}
			return m, nil

		case tea.KeyEnter:
			if len(m.items) > 0 && !m.showDetails {
				m.showDetails = true
			} else {
				m.showDetails = false
			}
			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(m.username), nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *InventoryModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Инвентарь"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	if len(m.items) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Инвентарь пуст"))
		return sb.String()
	}

	if m.showDetails {
		item := m.items[m.selected]
		sb.WriteString(ui.SelectedStyle.Render(item.Name))
		sb.WriteString("\n\n")
		sb.WriteString(ui.NormalStyle.Render("Количество: "))
		sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("%d", item.Quantity)))
		sb.WriteString("\n\n")
		sb.WriteString(ui.NormalStyle.Render("Описание:\n"))
		sb.WriteString(ui.NormalStyle.Render(item.ItemDescription))
	} else {
		for i, item := range m.items {
			if i == m.selected {
				sb.WriteString(ui.SelectedStyle.Render("> " + item.Name))
				sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf(" (x%d)", item.Quantity)))
			} else {
				sb.WriteString(ui.NormalStyle.Render("  " + item.Name))
				sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf(" (x%d)", item.Quantity)))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	if m.showDetails {
		sb.WriteString(ui.NormalStyle.Render("Enter - назад, Esc - в меню"))
	} else {
		sb.WriteString(ui.NormalStyle.Render("↑/↓ - выбор, Enter - подробности, Esc - в меню"))
	}

	return sb.String()
}
