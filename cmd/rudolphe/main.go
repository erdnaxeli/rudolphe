package main

import (
	"context"
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
	notifier := bot.NewNotifier(clock, repo)

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
	go notify(client, notifier)

	ctx := context.Background()
	log.Fatal(client.StartSync(ctx))
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
		sleep := doRefresh(client, refresher)
		slog.Info("Going to sleep before next refresh", "sleep", sleep)
		time.Sleep(sleep)
	}
}

func doRefresh(client matrix.Client, refresher bot.LimitedRefresher) time.Duration {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	result, sleep, err := refresher.Refresh(ctx)
	if err != nil {
		slog.Error("Unable to refresh leaderboards", "error", err)
	}

	for _, msg := range result.Messages {
		err := client.SendText(ctx, msg)
		if err != nil {
			slog.Error("Unable to send messages", "error", err)
		}
	}

	return sleep
}

func notify(client matrix.Client, notifier bot.Notifier) {
	sleep := notifier.GetSleep()
	slog.Info("Going to sleep before next notification", "sleep", sleep)
	time.Sleep(sleep)

	for {
		sleep := doNotify(client, notifier)
		slog.Info("Going to sleep before next notification", "sleep", sleep)
		time.Sleep(sleep)
	}
}

func doNotify(client matrix.Client, notifier bot.Notifier) time.Duration {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	result, sleep, err := notifier.GetNotification(ctx)
	if err != nil {
		slog.Error("Unabled to get notification", "error", err)
		return sleep
	}

	for _, msg := range result.Messages {
		err := client.SendText(ctx, msg)
		if err != nil {
			slog.Error("Unable to send notification messages", "error", err)
		}
	}

	if !result.LeaderBoard.IsEmpty() {
		err := client.SendLeaderBoard(ctx, result.LeaderBoard)
		if err != nil {
			slog.Error("Unable to send notification leaderboard", "error", err)
		}
	}

	return sleep
}
