package bot

import (
	"time"

	"github.com/erdnaxeli/rudolphe/aoc"
	"github.com/erdnaxeli/rudolphe/leaderboard"
)

type Result struct {
	LeaderBoard leaderboard.LeaderBoard
	Messages    []string
}

var EmptyResult = Result{}

type Bot interface {
	// Parse a message and return a response if it was a bot command.
	ParseCommand(msg string) (Result, error)
	// Refresh all leaderboards.
	//
	// If any update happened, it returns the corresponding messages.
	// It also returns the time to wait before the next refresh.
	Refresh() (Result, time.Duration, error)
}

const MIN_REFRESH_SLEEP = 15 * time.Minute

type bot struct {
	aocClient       aoc.Client
	leaderBoardRepo leaderboard.Repository

	lastRefresh time.Time
}

func New(
	aocClient aoc.Client,
	leaderBoardRepo leaderboard.Repository,
) bot {
	return bot{
		aocClient:       aocClient,
		leaderBoardRepo: leaderBoardRepo,
	}
}
