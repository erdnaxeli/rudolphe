package bot

import (
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
	Refresh() (Result, error)
}

type bot struct {
	aocClient       aoc.Client
	leaderBoardRepo leaderboard.Repository
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
