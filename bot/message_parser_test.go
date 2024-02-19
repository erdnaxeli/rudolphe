package bot_test

import (
	"context"
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

func (r *MRepo) GetLeaderBoard(ctx context.Context, year uint) (leaderboard.LeaderBoard, error) {
	args := r.Called(ctx, year)
	return args.Get(0).(leaderboard.LeaderBoard), args.Error(1)
}

func (r *MRepo) SaveLeaderBoard(ctx context.Context, year uint, leaderboard leaderboard.LeaderBoard) error {
	args := r.Called(ctx, year)
	return args.Error(0)
}

type MRefresher struct {
	mock.Mock
}

func (r *MRefresher) Refresh(ctx context.Context) (bot.Result, time.Duration, error) {
	args := r.Called(ctx)
	return args.Get(0).(bot.Result), args.Get(1).(time.Duration), args.Error(2)
}

func TestAutoRefreshMessageParser(t *testing.T) {
	// setup
	repo := &MRepo{}
	refresher := &MRefresher{}
	refresher.On("Refresh", mock.AnythingOfType("backgroundCtx")).Return(
		bot.Result{
			LeaderBoard: leaderboard.LeaderBoard{},
			Messages:    []string{"some message", "another one"},
		},
		15*time.Second,
		nil,
	)

	parser := bot.NewAutoRefreshMessageParser(refresher, repo)
	ctx := context.Background()

	// test
	result, err := parser.ParseMessage(ctx, "test")

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
