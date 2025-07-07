package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	authapi "lesta-start-battleship/cli/internal/api/auth"
	"lesta-start-battleship/cli/internal/cli/ui"
	"strings"
)

const ()

type OAuthModel struct {
	parent     tea.Model
	provider   string // "google" или "yandex"
	oauthURI   string
	deviceCode string // Для реального OAuth, здесь будет код устройства
	status     string // "waiting", "success", "error"
	username   string
	gold       int
	errorMsg   string
	authClient *authapi.Client
}

func NewOAuthModel(parent tea.Model, provider string, authClient *authapi.Client, oauthURL, deviceCode string) *OAuthModel {
	return &OAuthModel{
		parent:     parent,
		provider:   provider,
		status:     "waiting",
		oauthURI:   oauthURL,
		authClient: authClient,
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
			if m.status == "waiting" {
				// В реальности здесь будет API запрос
				// Сейчас просто ручной ввод результата
				return m, nil
			}
			// После успеха/ошибки Enter возвращает в родительское меню
			return m.parent, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	// Ручное управление результатом (для тестов)
	case OAuthTestSuccessMsg:
		m.status = "success"
		m.username = msg.Username
		m.gold = msg.Gold
		return m, nil

	case OAuthTestErrorMsg:
		m.status = "error"
		m.errorMsg = msg.Error
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
		sb.WriteString("1. Скопируйте ссыку:\n")
		sb.WriteString(ui.AlertStyle.Render(m.oauthURI))
		sb.WriteString("\n\n2. Откройте её в браузере\n\n")
		sb.WriteString("3. После авторизации нажмите Enter для проверки\n\n")
		sb.WriteString(ui.HelpStyle.Render("Enter - проверить, Esc - отмена"))

	case "success":
		sb.WriteString(ui.SuccessStyle.Render("Успешная авторизация!\n\n"))
		sb.WriteString(fmt.Sprintf("Игрок: %s\n", m.username))
		sb.WriteString(fmt.Sprintf("Золото: %d\n\n", m.gold))
		sb.WriteString(ui.HelpStyle.Render("Нажмите Enter чтобы продолжить"))

	case "error":
		sb.WriteString(ui.ErrorStyle.Render("Ошибка:\n"))
		sb.WriteString(m.errorMsg + "\n\n")
		sb.WriteString(ui.HelpStyle.Render("Нажмите Enter чтобы повторить"))
	}

	return sb.String()
}

// Тестовые сообщения для ручного управления
type OAuthTestSuccessMsg struct {
	Username string
	Gold     int
}

type OAuthTestErrorMsg struct {
	Error string
}
