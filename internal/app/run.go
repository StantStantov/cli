package app

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	cliModel "lesta-start-battleship/cli/internal/cli/initCli"
)

type App struct {
	program *tea.Program
}

func New(ctx context.Context) (*App, error) {
	initialModel := cliModel.NewCLI()

	program := tea.NewProgram(initialModel, tea.WithAltScreen())

	return &App{
		program: program,
	}, nil
}

func (a *App) Run() error {
	if _, err := a.program.Run(); err != nil {
		return err
	}
	return nil
}
