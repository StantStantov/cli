package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	guildStore "lesta-start-battleship/cli/store/guild"
	"strings"
)

const guildPerPage = 10

type GuildListModel struct {
	parent      tea.Model
	id          int
	username    string
	guilds      []guilds.GuildResponse
	currentPage int
	totalPages  int
	loading     bool
	errorMsg    string
	Clients     *clientdeps.Client
}

func NewGuildListModel(parent tea.Model, id int, username string, clients *clientdeps.Client) *GuildListModel {
	return &GuildListModel{
		parent:      parent,
		id:          id,
		username:    username,
		currentPage: 1,
		loading:     true,
		Clients:     clients,
	}
}

func (m *GuildListModel) Init() tea.Cmd {
	return m.loadGuilds
}

func (m *GuildListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			if m.currentPage > 1 {
				m.currentPage--
				return m, m.loadGuilds
			}
			return m, nil

		case tea.KeyRight:
			if m.currentPage < m.totalPages {
				m.currentPage++
				return m, m.loadGuilds
			}
			return m, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	case *guilds.GuildPagination:
		m.loading = false
		var activeGuilds []guilds.GuildResponse
		for _, guild := range msg.Items {
			if guild.IsActive {
				activeGuilds = append(activeGuilds, guild)
				guildStore.SetGuild(guild.Tag, guild)
			}
		}
		m.guilds = activeGuilds
		m.totalPages = msg.TotalPages
		if len(m.guilds) == 0 {
			m.errorMsg = "Нет активных гильдий"
		}
		return m, nil

	case error:
		m.loading = false
		m.errorMsg = msg.Error()
		return m, nil
	}

	return m, nil
}

func (m *GuildListModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Список гильдий"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("Страница %d/%d", m.currentPage, m.totalPages)))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString(ui.NormalStyle.Render("Загрузка списка гильдий..."))
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg + "\n"))
	}

	if len(m.guilds) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Нет активных гильдий"))
	} else {
		for _, guild := range m.guilds {
			line := fmt.Sprintf("%s - %s [%s]", guild.Title, guild.Description, guild.Tag)

			if guild.IsFull {
				sb.WriteString(ui.ErrorStyle.Render(line + " (Полная)\n"))
			} else {
				sb.WriteString(ui.NormalStyle.Render(line + "\n"))
			}
		}
	}

	sb.WriteString("\n")
	sb.WriteString(ui.HelpStyle.Render("←/→ - переключение страниц, Esc - назад"))

	return sb.String()
}

func (m *GuildListModel) loadGuilds() tea.Msg {
	ctx := context.Background()
	offset := (m.currentPage - 1) * guildPerPage
	guildsList, err := m.Clients.GuildsClient.GetGuilds(ctx, offset, guildPerPage)
	if err != nil {
		return err
	}
	return guildsList
}
