package aoc

import (
	"github.com/erdnaxeli/rudolphe/leaderboard"
)

const URL = "https://adventofcode.com"

type Client interface {
	GetLeaderBoard(year uint) leaderboard.LeaderBoard
}
