package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"LoLItemRecommender/internal/crawler"
	"LoLItemRecommender/internal/database"
	"LoLItemRecommender/internal/printer"
	"LoLItemRecommender/internal/queue"
	"LoLItemRecommender/internal/riotapi/gamedata"
)

var ErrNoAPIKey = errors.New("no riot api key set")

func getInitialPlayers(db *database.DB, gd *crawler.GameData) ([]*gamedata.Player, error) {
	players, err := db.GetLatestCrawledPlayers()
	if err != nil {
		return nil, err
	}
	if len(players) == 0 {
		return gd.InitWithChallengerPlayers()
	}
	excludedNames := make([]string, 0, len(players))
	for _, p := range players {
		excludedNames = append(excludedNames, p.SummonerName)
	}
	playersCrawled, err := db.GetAllPlayerNameCrawledExcept(excludedNames)
	if err != nil {
		return nil, err
	}
	gd.FlagPlayers(playersCrawled)
	return players, nil
}

func handleErrorsAndSignals(p *queue.Pool, ctx context.Context, cancel context.CancelFunc, signalChan chan os.Signal) {
	for {
		select {
		case e := <-p.GetErrorsChan():
			printer.Error("{-F_RED}Error received '%signalChan'", e.Error())
		case <-ctx.Done():
			printer.Debug("Stopped chan err")
			return
		case <-signalChan:
			printer.Info("Signal received, closing all")
			p.Close()
			cancel()
			return
		}
	}
}

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.CreateTables(); err != nil {
		log.Fatal(err)
	}
	if os.Getenv("RIOT_API_KEY") == "" {
		log.Fatal(ErrNoAPIKey)
	}
	gd, err := crawler.NewGameData("euw1", db)
	if err != nil {
		log.Fatal(err)
	}
	players, err := getInitialPlayers(db, gd)
	if err != nil {
		log.Fatal(err)
	}
	p := queue.NewPool(queue.CalculatePoolCap(players))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGQUIT)

	go handleErrorsAndSignals(p, ctx, cancel, signalChan)

	for _, player := range players {
		printer.Debug("Dispatch for player {-F_YELLOW}%s", player.SummonerName)
		p.Dispatch(func() error {
			return gd.CrawlPlayerData(player, p)
		})
	}
	p.WaitJobsToComplete()
	cancel()
	printer.Info("Jobs completed, channels closed")
}
