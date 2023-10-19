package bot

import (
	"github.com/erdnaxeli/rudolphe/leaderboard"
)

type Result struct {
	LeaderBoard leaderboard.LeaderBoard
	Messages    []string
}

var EmptyResult = Result{}
