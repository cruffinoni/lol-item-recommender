package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"LoLItemRecommender/internal/database"
	"LoLItemRecommender/internal/printer"
	"LoLItemRecommender/internal/riotapi/api"
	"LoLItemRecommender/internal/riotapi/gamedata"
	"LoLItemRecommender/internal/ui"
)

var ErrNoAPIKey = errors.New("no riot api key set")

func cancelUI(ctx context.Context, cancel context.CancelFunc, signalChan <-chan os.Signal) {
	select {
	case <-ctx.Done():
		printer.Debug("Stopped chan err")
		return
	case <-signalChan:
		printer.Info("Signal received, closing all")
		cancel()
		return
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	client := api.NewClient()
	em := api.NewEndpointsManager(os.Getenv("RIOT_API_KEY"), "euw1")
	staticData := gamedata.NewStaticData(em, client)
	if err := staticData.RetrieveAPIVersion(); err != nil {
		log.Fatal(err)
	}
	if err := staticData.RetrieveChampionsStats(); err != nil {
		log.Fatal(err)
	}
	if err := staticData.RetrieveItems(); err != nil {
		log.Fatal(err)
	}

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGQUIT)
	go cancelUI(ctx, cancel, signalChan)
	c := ui.NewConsole(ctx)
	c.AskUserChampions(staticData, db)
}
