package ui

import (
	"bufio"
	"context"
	"errors"
	"os"
	"strings"

	"LoLItemRecommender/internal/database"
	"LoLItemRecommender/internal/printer"
	"LoLItemRecommender/internal/riotapi/gamedata"
	"LoLItemRecommender/internal/style"
)

type Console struct {
	blueTeam []*gamedata.ChampionStats
	redTeam  []*gamedata.ChampionStats
	reader   *bufio.Reader
	quit     chan struct{}
}

const (
	BlueTeam = "blue"
	RedTeam  = "red"
)

func NewConsole(ctx context.Context) *Console {
	c := &Console{
		blueTeam: make([]*gamedata.ChampionStats, 0),
		redTeam:  make([]*gamedata.ChampionStats, 0),
		reader:   bufio.NewReader(os.Stdin),
		quit:     make(chan struct{}),
	}
	go func(ctx context.Context, c *Console) {
		select {
		case <-ctx.Done():
			close(c.quit)
		}
	}(ctx, c)
	return c
}

func (c *Console) DisplayInstructions() {
	printer.Printf("{-F_CYAN,BOLD}════════════════════════════════════════════════")
	printer.Printf("{-F_CYAN,BOLD}       ItemResponse Advisor - League of Legends       ")
	printer.Printf("{-F_CYAN,BOLD}════════════════════════════════════════════════")

	printer.Print("Follow the instructions to get item suggestions for Samira based on the game composition")
	printer.Print("--------------------------------------------------------")
}

func (c *Console) AddChampionToTeam(team string, champion *gamedata.ChampionStats) {
	if team == BlueTeam {
		c.blueTeam = append(c.blueTeam, champion)
	} else {
		c.redTeam = append(c.redTeam, champion)
	}
}

func (c *Console) readNonBlockingInput() <-chan string {
	inputChan := make(chan string)
	go func() {
		input, err := c.reader.ReadString('\n')
		if err != nil {
			printer.PrintError(err)
		}
		inputChan <- strings.TrimSpace(input)
	}()
	return inputChan
}

var ErrContextCanceled = errors.New("context canceled")

func (c *Console) AskForChampionTeam(champion *gamedata.ChampionStats) error {
	for {
		printer.Print("In which team the champion is ? {-F_BLUE,BOLD}Blue team{-RESET} or {-F_RED,BOLD}Red team {-RESET}?")
		var team string
		select {
		case team = <-c.readNonBlockingInput():
			switch team {
			case BlueTeam:
				printer.Debug("Adding %v to blue team", champion.Name)
				c.blueTeam = append(c.blueTeam, champion)
				return nil
			case RedTeam:
				printer.Debug("Adding %v to red team", champion.Name)
				c.redTeam = append(c.redTeam, champion)
				return nil
			default:
				printer.Error("The team doesn't not exist")
			}
		case <-c.quit:
			printer.Debug("Channel quit called from ask for champ")
			return ErrContextCanceled
		}
	}
}
func (c *Console) Display() {
	printer.Printf("{-F_BLUE,BOLD}Blue Team{-RESET} composition:")
	for i, champion := range c.blueTeam {
		printer.Printf("{-F_BLUE}[%d] %s (id %v)", i+1, champion.Name, champion.Key)
	}

	printer.Printf("\n{-F_RED,BOLD}Red Team{-RESET} composition:")
	for i, champion := range c.redTeam {
		printer.Printf("{-F_RED}[%d] %s (id %v)", i+1, champion.Name, champion.Key)
	}
}

func (c *Console) AskUserChampions(sd *gamedata.StaticData, db *database.DB) {
	c.DisplayInstructions()
	for {
		printer.Print("Please enter the names of the champions for which you want item suggestions. Type {-BOLD}'end'{-RESET} to finish entering.")
		printer.Print("Enter a champion name: ")

		select {
		case <-c.quit:
			printer.Debug("Channel quit called")
			return
		case input := <-c.readNonBlockingInput():
			if input == "" {
				continue
			}

			input = strings.ToLower(input)

			switch input {
			case "end":
				return
			case "clear":
				c.redTeam = make([]*gamedata.ChampionStats, 0)
				c.blueTeam = make([]*gamedata.ChampionStats, 0)
				continue
			case "display", "show":
				c.Display()
				continue
			case "search":
				if len(c.blueTeam) == 0 || len(c.redTeam) == 0 {
					printer.Error("not enough participant to make a search")
					continue
				}
				p, err := db.GetMatchesWithChampions(c.blueTeam, c.redTeam)
				if err != nil {
					printer.PrintError(err)
					return
				}
				printer.Debug("%v participants found like %d matches", len(p), len(p)/10)
				continue
			}

			input = style.ToTitleCase(input)

			if input == "Wukong" { // Riot API call the champion Wukong: "MonkeyKing"
				input = "MonkeyKing"
			}
			s := sd.GetChampionStats(input)
			if s == nil {
				printer.Error("'%s' doesn't exist", input)
				continue
			}
			printer.Printf("Champion '%s' found", s.Name)
			if err := c.AskForChampionTeam(s); err != nil && err == ErrContextCanceled {
				return
			}
		}
	}
}
