package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/erdnaxeli/rudolphe/aoc"
	"github.com/erdnaxeli/rudolphe/bot"
	"github.com/erdnaxeli/rudolphe/leaderboard"
	"github.com/erdnaxeli/rudolphe/matrix"
)

type Config struct {
	LeaderBoardID string
	SessionCookie string
	HomeServerURL string
	UserID        string
	AccessToken   string
	RoomID        string
}

func main() {
	config := getConfig()
	repo, err := leaderboard.NewSqliteRepository("db.sqlite")
	if err != nil {
		log.Fatalf("Unable to create repository: %v", err)
	}

	aocClient, err := aoc.NewJsonClient(
		config.SessionCookie,
		config.LeaderBoardID,
	)
	if err != nil {
		log.Fatalf("Error while creating AOC client: %v", err)
	}

	clock := bot.DefaultClock{}
	refresher := bot.NewRefresher(aocClient, clock, repo)
	minRefresher := bot.NewMinLimitedRefresher(clock, refresher)
	monthRefresher := bot.NewMonthLimiterRefresher(clock, refresher)
	messageParser := bot.NewAutoRefreshMessageParser(minRefresher, repo)

	client, err := matrix.NewClient(
		matrix.Config{
			HomeserverURL: config.HomeServerURL,
			UserID:        config.UserID,
			AccessToken:   config.AccessToken,
			RoomID:        config.RoomID,
			MessageParser: messageParser,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	go refresh(client, monthRefresher)
	log.Fatal(client.StartSync())
}

func getConfig() Config {
	content, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error while reading config: %v", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Error while parsing config: %v", err)
	}

	return config
}

func refresh(client matrix.Client, refresher bot.LimitedRefresher) {
	for {
		result, sleep, err := refresher.Refresh()
		if err != nil {
			slog.Error("Unable to refresh leaderboards", "error", err)
		}

		for _, msg := range result.Messages {
			err := client.SendText(msg)
			if err != nil {
				slog.Error("Unable to send messages", "error", err)
			}
		}

		slog.Info("Going to sleep before next refresh", "sleep", sleep)
		time.Sleep(sleep)
	}
}
