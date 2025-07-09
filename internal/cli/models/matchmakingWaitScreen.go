package models

import (
	"fmt"
	"lesta-battleship/cli/internal/cli/ui"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type MatchmakingWaitScreenModel struct {
	username string

	startTime time.Time
	endTime   time.Time
}

func NewMatchmakingWaitScreenModel(username string) *MatchmakingWaitScreenModel {
	now := time.Now()

	return &MatchmakingWaitScreenModel{
		username: username,

		startTime: now,
		endTime: now,
	}
}

func (m *MatchmakingWaitScreenModel) Init() tea.Cmd {
	return tickEvery()
}

func (m *MatchmakingWaitScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewMatchmakingModel(m.username), nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tickMsg:
		m.endTime = time.Time(msg)
		return m, tickEvery()
	}

	return m, nil
}

func (m *MatchmakingWaitScreenModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	fmt.Fprintf(&sb, "Время прошло: %s", m.endTime.Sub(m.startTime).Round(time.Second))

	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Esc - выход"))

	return sb.String()
}

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
