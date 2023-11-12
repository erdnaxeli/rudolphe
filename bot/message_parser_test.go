package bot_test

import (
	"testing"
	"time"

	"github.com/erdnaxeli/rudolphe/bot"
	"github.com/erdnaxeli/rudolphe/leaderboard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MRepo struct {
	mock.Mock
}

func (r *MRepo) GetLeaderBoard(year uint) (leaderboard.LeaderBoard, error) {
	args := r.Called(year)
	return args.Get(0).(leaderboard.LeaderBoard), args.Error(1)
}

func (r *MRepo) SaveLeaderBoard(year uint, leaderboard leaderboard.LeaderBoard) error {
	args := r.Called(year)
	return args.Error(0)
}

type MRefresher struct {
	mock.Mock
}

func (r *MRefresher) Refresh() (bot.Result, time.Duration, error) {
	args := r.Called()
	return args.Get(0).(bot.Result), args.Get(1).(time.Duration), args.Error(2)
}

func TestAutoRefreshMessageParser(t *testing.T) {
	// setup
	repo := &MRepo{}
	refresher := &MRefresher{}
	refresher.On("Refresh").Return(
		bot.Result{
			LeaderBoard: leaderboard.LeaderBoard{},
			Messages:    []string{"some message", "another one"},
		},
		15*time.Second,
		nil,
	)

	parser := bot.NewAutoRefreshMessageParser(refresher, repo)

	// test
	result, err := parser.ParseMessage("test")

	// assertions
	assert.ErrorIs(t, err, bot.ErrUnknownCommand)
	assert.Equal(t,
		bot.Result{
			LeaderBoard: leaderboard.LeaderBoard{},
			Messages:    []string{"some message", "another one"},
		},
		result,
	)
}
