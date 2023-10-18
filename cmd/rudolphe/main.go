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
	repo, err := leaderboard.NewSqliteRepository("db.sqlite")
	if err != nil {
		log.Fatalf("Unable to create repository: %v", err)
	}

	content, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error while reading config: %v", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Error while parsing config: %v", err)
	}

	aocClient, err := aoc.NewJsonClient(
		config.SessionCookie,
		config.LeaderBoardID,
	)
	if err != nil {
		log.Fatalf("Error while creating AOC client: %v", err)
	}

	bot := bot.New(aocClient, repo)
	client, err := matrix.NewClient(
		matrix.Config{
			HomeserverURL: config.HomeServerURL,
			UserID:        config.UserID,
			AccessToken:   config.AccessToken,
			RoomID:        config.RoomID,
			Bot:           bot,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = startUpdater(client, bot)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(client.StartSync())
	//mClient.SendText()
}

func startUpdater(client matrix.Client, bot bot.Bot) error {
	go func() {
		for {
			now := time.Now()
			result, err := bot.Refresh()
			if err != nil {
				slog.Error("Unable to refresh leaderboards", "error", err)
			}

			for _, msg := range result.Messages {
				err := client.SendText(msg)
				if err != nil {
					slog.Error("Unable to send messages", "error", err)
				}
			}

			var sleep time.Duration
			if time.January <= now.Month() && now.Month() <= time.November {
				sleep = 5 * time.Hour
			} else {
				sleep = 15 * time.Minute
			}

			slog.Info("Going to sleep before next refresh", "sleep", sleep)
			time.Sleep(sleep)
		}
	}()

	return nil
}
