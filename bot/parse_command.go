package bot

import (
	"strconv"
	"time"

	"golang.org/x/exp/slog"
)

func (b bot) ParseCommand(msg string) (Result, error) {
	result, err := b.parseCommand(msg)
	if err != nil {
		return result, err
	}

	if time.Since(b.lastRefresh) >= 15*time.Minute {
		refreshResult, err := b.Refresh()
		if err != nil {
			slog.Warn("Error while auto refreshing", "error", err)
		} else {
			result.Messages = append(result.Messages, refreshResult.Messages...)
		}
	}

	return result, nil
}

func (b bot) parseCommand(msg string) (Result, error) {
	if msg == "!lb" {
		now := time.Now()
		year := now.Year()

		if now.Month() < 12 {
			year--
		}

		lb, err := b.leaderBoardRepo.GetLeaderBoard(uint(year))
		if err != nil {
			return EmptyResult, err
		}

		return Result{LeaderBoard: lb}, nil
	} else if len(msg) == 8 && msg[:4] == "!lb " {
		year, err := strconv.Atoi(msg[4:8])
		if err != nil {
			return Result{Messages: []string{"Année invalide"}}, nil
		}

		lb, err := b.leaderBoardRepo.GetLeaderBoard(uint(year))
		if err != nil {
			return EmptyResult, err
		}

		return Result{LeaderBoard: lb}, nil
	}

	return EmptyResult, ErrUnknownCommand
}
