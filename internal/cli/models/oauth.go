package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

const (
	expiresInSeconds = 600
	intervalSeconds  = 5
)

type OAuthModel struct {
	parent     tea.Model
	provider   string // "google" или "yandex"
	oauthURI   string
	deviceCode string // Для реального OAuth, здесь будет код устройства
	status     string // "waiting", "success", "error"
	id         int
	username   string
	gold       int
	errorMsg   string
	Clients    *clientdeps.Client
}

func NewOAuthModel(parent tea.Model, provider string, client *clientdeps.Client, oauthURL, deviceCode string) *OAuthModel {
	return &OAuthModel{
		parent:     parent,
		provider:   provider,
		status:     "waiting",
		oauthURI:   oauthURL,
		deviceCode: deviceCode,
		Clients:    client,
	}
}

func (m *OAuthModel) Init() tea.Cmd {
	return nil // Никаких автоматических действий
}

func (m *OAuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			/*switch m.status {
			case "waiting", "error":
				m.status = "pending"
				m.errorMsg = ""
				return m, m.pollingOAuth()
			case "success":
				return NewMainMenuModel(m.id, m.username, m.gold), nil
			case "pending":
				return m, nil
			}*/
			if m.status == "waiting" {
				// В реальности здесь будет API запрос
				// Сейчас просто ручной ввод
				//m.status = "pending"
				m.errorMsg = ""

				return m, func() tea.Msg { return OAuthPollingResultMsg{Error: "Пользователь не найден"} } // Здесь будет реальный результат
			} else if m.status == "error" {
				return NewMainMenuModel(m.id, m.username, m.gold, m.Clients), nil
			}
			// После успеха/ошибки Enter возвращает в родительское меню
			return m.parent, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	case OAuthPollingResultMsg:
		if msg.Error != "" {
			m.status = "error"
			m.errorMsg = msg.Error
			return m, nil
		}
		m.status = "success"
		m.id = msg.ID
		m.username = msg.Username
		m.gold = msg.Gold
		return m, nil
	}

	return m, nil
}

func (m *OAuthModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render(fmt.Sprintf("Авторизация через %s", strings.Title(m.provider))))
	sb.WriteString("\n\n")

	switch m.status {
	case "waiting":
		sb.WriteString("1. Скопируйте ссылку:\n")
		sb.WriteString(ui.AlertStyle.Render(m.oauthURI))
		sb.WriteString("\n\n2. Откройте её в браузере\n\n")
		sb.WriteString(fmt.Sprintf("3. Введите код устройства: %s\n\n", m.deviceCode))
		sb.WriteString("3. После авторизации нажмите Enter для проверки\n\n")
		sb.WriteString(ui.HelpStyle.Render("Enter - проверить, Esc - отмена"))

	case "pending":
		sb.WriteString("Ожидание подтверждения авторизации...\n")
		sb.WriteString(ui.HelpStyle.Render("Esc - отмены"))

	case "success":
		sb.WriteString(ui.SuccessStyle.Render("Успешная авторизация!\n\n"))
		sb.WriteString(fmt.Sprintf("Игрок: %s\n", m.username))
		sb.WriteString(fmt.Sprintf("Золото: %d\n\n", m.gold))
		sb.WriteString(ui.HelpStyle.Render("Нажмите Enter чтобы продолжить"))

	case "error":
		sb.WriteString(ui.ErrorStyle.Render(fmt.Sprintf("Ошибка: %s\n\n", m.errorMsg)))
		sb.WriteString(ui.HelpStyle.Render("Enter - повторить, Esc - назад"))
	}

	return sb.String()
}

func (m *OAuthModel) pollingOAuth() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		_, profile, err := m.Clients.AuthClient.CompleteOAuthPolling(ctx, m.provider, m.deviceCode, expiresInSeconds, intervalSeconds)
		if err != nil {
			return OAuthPollingResultMsg{Error: err.Error()}
		}
		return OAuthPollingResultMsg{ID: profile.ID, Username: profile.Username, Gold: profile.Currency.Gold}
	}
}
