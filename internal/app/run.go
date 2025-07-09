package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/auth"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/api/inventory"
	"lesta-start-battleship/cli/internal/api/scoreboard"
	"lesta-start-battleship/cli/internal/api/shop"
	cliModel "lesta-start-battleship/cli/internal/cli/initCli"
	"lesta-start-battleship/cli/internal/clientdeps"
)

const (
	authURL       = "https://battleship-lesta-start.ru/"
	guildsURL     = "https://battleship-lesta-start.ru/guild/"
	inventoryURL  = "https://battleship-lesta-start.ru/inventory/"
	scoreboardURL = "https://battleship-lesta-start.ru/scoreboard/"
	shopURL       = "https://battleship-lesta-start.ru/shop/"
)

type App struct {
	program *tea.Program
}

func New() (*App, error) {
	initialClients, err := initClients()
	if err != nil {
		return nil, err
	}

	initialModel := cliModel.NewCLI(initialClients)

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

func initClients() (*clientdeps.Client, error) {
	authClient, err := auth.NewClient(authURL)
	if err != nil {
		return nil, err
	}

	guildsClient, err := guilds.NewClient(guildsURL)
	if err != nil {
		return nil, err
	}

	inventoryClient, err := inventory.NewClient(inventoryURL)
	if err != nil {
		return nil, err
	}

	scoreboardClient := scoreboard.NewClient(scoreboardURL, nil)
	shopClient := shop.NewClient(shopURL)

	return &clientdeps.Client{
		AuthClient:       authClient,
		GuildsClient:     guildsClient,
		InventoryClient:  inventoryClient,
		ScoreboardClient: scoreboardClient,
		ShopClient:       shopClient,
	}, nil
}
