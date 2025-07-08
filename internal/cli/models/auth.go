package models

import (
	"github.com/charmbracelet/bubbletea"
	"lesta-battleship/cli/internal/cli/ui"
	"strings"
)

type AuthModel struct {
	login       string
	password    string
	activeField int
	errorMsg    string
	registering bool
}

func NewAuthModel() *AuthModel {
	return &AuthModel{
		activeField: 0,
	}
}

func (m *AuthModel) Init() tea.Cmd {
	return nil
}

func (m *AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.activeField == 1 {
				m.errorMsg = ""
				//m.Error = "Неверное имя или пароль"
				if m.login == "" || m.password == "" {
					m.errorMsg = "Введите логин и пароль"
					return m, nil
				}
				if m.registering {
					//ручка с регистрацией
					m.login = ""
					m.password = ""
					m.registering = false
					return m, nil
				}
				//ручка с авторизацией
				if m.errorMsg != "" {
					return m, nil
				} else {
					return m, func() tea.Msg {
						return AuthSuccessMsg{
							Token:    "dummy_token_" + m.login,
							Username: m.login,
						}
					}
				}
			}
			m.activeField = 1
			return m, nil

		case tea.KeyTab:
			m.activeField = (m.activeField + 1) % 2
			return m, nil

		case tea.KeyBackspace:
			if m.activeField == 0 && len(m.login) > 0 {
				m.login = m.login[:len(m.login)-1]
			} else if m.activeField == 1 && len(m.password) > 0 {
				m.password = m.password[:len(m.password)-1]
			}
			return m, nil

		case tea.KeyRunes:
			if m.activeField == 0 {
				m.login += string(msg.Runes)
			} else {
				m.password += string(msg.Runes)
			}
			return m, nil

		case tea.KeyCtrlN:
			m.registering = !m.registering
			return m, nil

		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *AuthModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")

	if m.registering {
		sb.WriteString("Регистрация\n\n")
	} else {
		sb.WriteString("Авторизация\n\n")
	}

	loginLabel := "Логин: "
	if m.activeField == 0 {
		sb.WriteString(ui.SelectedStyle.Render(loginLabel))
	} else {
		sb.WriteString(ui.NormalStyle.Render(loginLabel))
	}
	sb.WriteString(m.login)
	if m.activeField == 0 {
		sb.WriteString("_")
	}
	sb.WriteString("\n")

	passwordLabel := "Пароль: "
	if m.activeField == 1 {
		sb.WriteString(ui.SelectedStyle.Render(passwordLabel))
	} else {
		sb.WriteString(ui.NormalStyle.Render(passwordLabel))
	}
	sb.WriteString(strings.Repeat("*", len(m.password)))
	if m.activeField == 1 {
		sb.WriteString("_")
	}
	sb.WriteString("\n")

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg + "\n"))
	}

	sb.WriteString(ui.NormalStyle.Render("\nTab - переключение полей, Enter - подтвердить"))
	sb.WriteString(ui.NormalStyle.Render("\nCtrl+N - переключить регистрацию/авторизацию, Esc/Ctrl+C - выход\n"))

	return sb.String()
}
