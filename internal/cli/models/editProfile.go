package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/cli/ui"
	"strings"
)

type EditProfileModel struct {
	username    string
	tempNick    string
	tempPass    string
	activeTab   int // 0 - ник, 1 - пароль
	activeField int
	errorMsg    string
	successMsg  string
}

func NewEditProfileModel(username string) *EditProfileModel {
	return &EditProfileModel{
		username:    username,
		tempNick:    username,
		activeTab:   0,
		activeField: 0,
	}
}

func (m *EditProfileModel) Init() tea.Cmd {
	return nil
}

func (m *EditProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.activeTab == 0 {
				m.activeTab = 1
			} else {
				m.activeTab = 0
			}
			m.activeField = 0
			return m, nil

		case tea.KeyEnter:
			if m.activeTab == 0 && m.activeField == 1 {
				if len(m.tempNick) < 3 {
					m.errorMsg = "Ник должен быть не менее 3 символов"
				} else {
					m.username = m.tempNick
					m.errorMsg = ""
					m.tempNick = ""
					m.successMsg = "Ник успешно изменен!"
					return m, func() tea.Msg {
						return UsernameChangeMsg{NewUsername: m.tempNick}
					}
				}
			} else if m.activeTab == 1 && m.activeField == 1 {
				if len(m.tempPass) < 6 {
					m.errorMsg = "Пароль должен быть не менее 6 символов"
				} else {
					m.errorMsg = ""
					m.successMsg = "Пароль успешно изменен!"
					m.tempPass = ""
				}
			} else {
				m.activeField = 1
			}
			return m, nil

		case tea.KeyBackspace:
			if m.activeTab == 0 && m.activeField == 0 && len(m.tempNick) > 0 {
				m.tempNick = m.tempNick[:len(m.tempNick)-1]
			} else if m.activeTab == 1 && m.activeField == 0 && len(m.tempPass) > 0 {
				m.tempPass = m.tempPass[:len(m.tempPass)-1]
			}
			return m, nil

		case tea.KeyRunes:
			if m.activeTab == 0 && m.activeField == 0 {
				m.tempNick += string(msg.Runes)
			} else if m.activeTab == 1 && m.activeField == 0 {
				m.tempPass += string(msg.Runes)
			}
			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(m.username), nil
		}
	}
	return m, nil
}

func (m *EditProfileModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Редактирование профиля"))
	sb.WriteString("\n\n")

	nickTab := "Ник"
	passTab := "Пароль"
	if m.activeTab == 0 {
		nickTab = ui.SelectedStyle.Render(nickTab)
	} else {
		nickTab = ui.NormalStyle.Render(nickTab)
	}
	if m.activeTab == 1 {
		passTab = ui.SelectedStyle.Render(passTab)
	} else {
		passTab = ui.NormalStyle.Render(passTab)
	}
	sb.WriteString(nickTab + " | " + passTab)
	sb.WriteString("\n\n")

	if m.activeTab == 0 {
		sb.WriteString("Текущий ник: " + m.username + "\n\n")
		sb.WriteString("Новый ник:\n")
		if m.activeField == 0 {
			sb.WriteString(ui.SelectedStyle.Render("> " + m.tempNick + "_"))
		} else {
			sb.WriteString(" " + m.tempNick)
			sb.WriteString("\n\n")
			sb.WriteString(ui.SuccessStyle.Render("Нажмите Enter для сохранения"))
		}
	} else {
		sb.WriteString("Новый пароль:\n")
		if m.activeField == 0 {
			sb.WriteString(ui.SelectedStyle.Render("> " + strings.Repeat("*", len(m.tempPass)) + "_"))
		} else {
			sb.WriteString(" " + strings.Repeat("*", len(m.tempPass)))
			sb.WriteString("\n\n")
			sb.WriteString(ui.SuccessStyle.Render("Нажмите Enter для сохранения"))
		}
	}

	if m.errorMsg != "" {
		sb.WriteString("\n\n")
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg))
	}
	if m.successMsg != "" {
		sb.WriteString("\n\n")
		sb.WriteString(ui.SuccessStyle.Render(m.successMsg))
	}

	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Tab - переключение вкладок, Enter - подтвердить, Esc - назад"))

	return sb.String()
}
