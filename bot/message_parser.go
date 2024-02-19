package bot

import (
	"context"
	"strconv"
	"time"

	"github.com/erdnaxeli/rudolphe/leaderboard"
	"golang.org/x/exp/slog"
)

type MessageParser interface {
	// Parse a message and return a response if it was a bot command.
	ParseMessage(ctx context.Context, msg string) (Result, error)
}

// AutoRefreshMessageParser parse messages and trigger a refresh.
//
// It first refresh the data, then parse the message. Errors happening during
// the refresh are logged but not returned.
type AutoRefreshMessageParser struct {
	refresher LimitedRefresher
	repo      leaderboard.Repository
}

func NewAutoRefreshMessageParser(
	refresher LimitedRefresher,
	repo leaderboard.Repository,
) AutoRefreshMessageParser {
	return AutoRefreshMessageParser{
		refresher: refresher,
		repo:      repo,
	}
}

func (m AutoRefreshMessageParser) ParseMessage(ctx context.Context, msg string) (Result, error) {
	refreshResult, _, err := m.refresher.Refresh(ctx)
	if err != nil {
		slog.Warn("Error while auto refreshing", "error", err)
	}

	result, err := m.doParseMessage(ctx, msg)
	result.Messages = append(refreshResult.Messages, result.Messages...)
	return result, err
}

func (m AutoRefreshMessageParser) doParseMessage(ctx context.Context, msg string) (Result, error) {
	year := -1

	if msg == "!lb" {
		now := time.Now()
		year = now.Year()

		if now.Month() < 12 {
			year--
		}
	} else if len(msg) == 8 && msg[:4] == "!lb " {
		var err error
		year, err = strconv.Atoi(msg[4:8])
		if err != nil {
			return Result{Messages: []string{"Année invalide"}}, nil
		}
	}

	if year > 0 {
		lb, err := m.repo.GetLeaderBoard(ctx, uint(year))
		if err != nil {
			return EmptyResult, err
		}

		return Result{LeaderBoard: lb}, nil
	}

	return EmptyResult, ErrUnknownCommand
}
