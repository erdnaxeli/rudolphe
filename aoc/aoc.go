package aoc

import (
	"context"

	"github.com/erdnaxeli/rudolphe/leaderboard"
)

const URL = "https://adventofcode.com"

type Client interface {
	GetLeaderBoard(ctx context.Context, year uint) (leaderboard.LeaderBoard, error)
}
