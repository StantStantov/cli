package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/cli/handlers"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

const (
	pageSize = 5
)

type ScoreboardModel struct {
	username     string
	gold         int
	activeTab    int // 0-моя, 1-игроки, 2-гильдия
	myStats      handlers.PlayerStats
	playersStats []handlers.PlayerStats
	guildStats   handlers.GuildInternalStats
	err          error
	currentPage  int
	totalPages   int
	tableWidth   int
	Clients      *clientdeps.Client
}

func NewScoreboardModel(username string, gold int, myStats handlers.PlayerStats, clients *clientdeps.Client) *ScoreboardModel {
	return &ScoreboardModel{
		username:    username,
		gold:        gold,
		activeTab:   0,
		myStats:     myStats,
		currentPage: 1,
		tableWidth:  80,
		Clients:     clients,
	}
}

func (m *ScoreboardModel) Init() tea.Cmd {
	return m.loadStats
}

func (m *ScoreboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.tableWidth = msg.Width - 10
		return m, nil

	case handlers.PlayerStats:
		m.myStats = msg
		return m, nil

	case []handlers.PlayerStats:
		m.playersStats = msg
		m.totalPages = (len(msg) + pageSize - 1) / pageSize
		return m, nil

	case handlers.GuildInternalStats:
		m.guildStats = msg
		m.totalPages = (len(msg.Members) + pageSize - 1) / pageSize
		return m, nil

	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			m.activeTab = (m.activeTab - 1 + 3) % 3
			m.currentPage = 1
			return m, m.loadStats

		case tea.KeyRight:
			m.activeTab = (m.activeTab + 1) % 3
			m.currentPage = 1
			return m, m.loadStats

		case tea.KeyDown:
			if m.currentPage < m.totalPages {
				m.currentPage++
			}
			return m, nil

		case tea.KeyUp:
			if m.currentPage > 1 {
				m.currentPage--
			}
			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(m.username, m.gold, m.Clients), nil
		}
	}
	return m, nil
}

func (m *ScoreboardModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Рейтинги"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	tabs := []string{"Моя статистика", "Игроки", "Моя гильдия"}
	for i, tab := range tabs {
		if i == m.activeTab {
			sb.WriteString(ui.SelectedStyle.Render(" [" + tab + "] "))
		} else {
			sb.WriteString(ui.NormalStyle.Render(tab + " "))
		}
	}
	sb.WriteString("\n\n")

	if m.err != nil {
		sb.WriteString(ui.ErrorStyle.Render("Ошибка: " + m.err.Error()))
		return sb.String()
	}

	switch m.activeTab {
	case 0:
		sb.WriteString(fmt.Sprintf("Игрок: %s\n", m.myStats.Username))
		sb.WriteString(fmt.Sprintf("Ранг: %d\n", m.myStats.Rank))
		sb.WriteString(fmt.Sprintf("Победы: %d\n", m.myStats.Wins))
		sb.WriteString(fmt.Sprintf("Поражения: %d\n", m.myStats.Losses))
		sb.WriteString(fmt.Sprintf("Всего побед: %d\n", m.myStats.TotalGames))
		winRate := float64(m.myStats.Wins) / float64(m.myStats.TotalGames) * 100
		sb.WriteString(fmt.Sprintf("Процент побед: %.1f%%", winRate))

	case 1:
		start := (m.currentPage - 1) * pageSize
		end := start + pageSize
		if end > len(m.playersStats) {
			end = len(m.playersStats)
		}
		sb.WriteString(m.renderPlayersTable(m.playersStats[start:end]))

	case 2:
		start := (m.currentPage - 1) * pageSize
		end := start + pageSize
		if end > len(m.guildStats.Members) {
			end = len(m.guildStats.Members)
		}
		sb.WriteString(m.renderGuildStats())
		sb.WriteString(m.renderGuildMembersTable(m.guildStats.Members[start:end]))
	}

	if m.activeTab > 0 && m.totalPages > 1 {
		sb.WriteString(fmt.Sprintf("\nСтраница %d/%d", m.currentPage, m.totalPages))
	}

	sb.WriteString("\n\n")
	helpText := "←/→ - вкладки"
	if m.activeTab > 0 {
		helpText += ", PgUp/PgDown - страницы"
	}
	sb.WriteString(ui.NormalStyle.Render(helpText + ", Esc - назад"))

	return sb.String()
}

func (m *ScoreboardModel) renderPlayersTable(players []handlers.PlayerStats) string {
	headers := []string{"Ранг", "Игрок", "Победы", "Поражения", "% побед"}
	widths := []int{6, 20, 8, 10, 10}

	table := ui.NewTable(m.tableWidth, widths)
	table.AddHeader(headers)

	for _, p := range players {
		winRate := float64(p.Wins) / float64(p.TotalGames) * 100
		table.AddRow([]string{
			fmt.Sprintf("%d", p.Rank),
			p.Username,
			fmt.Sprintf("%d", p.Wins),
			fmt.Sprintf("%d", p.Losses),
			fmt.Sprintf("%.1f%%", winRate),
		})
	}

	return table.Render()
}

func (m *ScoreboardModel) renderGuildStats() string {
	return fmt.Sprintf(
		"Гильдия: [%s] %s\nОчки ярости: %d\n\n",
		m.guildStats.GuildTag, m.guildStats.GuildName, m.guildStats.TotalRage,
	)
}

func (m *ScoreboardModel) renderGuildMembersTable(members []handlers.GuildsMemberStats) string {
	headers := []string{"Игрок", "Победы", "Сундуки", "Ярость", "Войны", "Вклад"}
	widths := []int{20, 8, 10, 8, 8, 10}

	table := ui.NewTable(m.tableWidth, widths)
	table.AddHeader(headers)

	for _, p := range members {
		table.AddRow([]string{
			p.Username,
			fmt.Sprintf("%d", p.Wins),
			fmt.Sprintf("%d", p.GuildChests),
			fmt.Sprintf("%d", p.GuildRagePoints),
			fmt.Sprintf("%d", p.WarWins),
			fmt.Sprintf("%d", p.TotalDonations),
		})
	}

	return table.Render()
}

func (m *ScoreboardModel) loadStats() tea.Msg {
	token := "dummy_token_" + m.username

	switch m.activeTab {
	case 0:
		stats, err := handlers.MyStatsHandler(token)
		if err != nil {
			return err
		}
		return stats
	case 1:
		stats, err := handlers.PlayersStatsHandler()
		if err != nil {
			return err
		}
		return stats
	case 2:
		stats, err := handlers.GuildInternalStatsHandler(token)
		if err != nil {
			return err
		}
		return stats
	}
	return nil
}
