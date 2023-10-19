package bot

import (
	"fmt"
	"time"

	"github.com/erdnaxeli/rudolphe/leaderboard"
)

// Notifier is used to send a notification to users.
type Notifier interface {
	// GetNotification returns the notification and how long to wait before.
	//
	// The notification is a Result object. It can contains a leaderboard and / or
	// messages.
	// If any error happens the result should be discarded but the duration to wait
	// will still be the correct one.
	GetNotification() (Result, time.Duration, error)
	// GetSleep returns the time to wait until the next notification.
	//
	// It is useful on start to just wait until the next notification, and then start
	// the cycle with GetNotification().
	GetSleep() time.Duration
}

const (
	AOC_URL                 = "https://adventofcode.com"
	AOC_LOCATION_UTC_OFFSET = -5
	NEW_PUZZLE_FMT          = "Nouveau puzzle : %s/%d/day/%d"
)

type DefaultNotifier struct {
	clock Clock
	repo  leaderboard.Repository
}

func NewNotifier(clock Clock, repo leaderboard.Repository) DefaultNotifier {
	return DefaultNotifier{
		clock: clock,
		repo:  repo,
	}
}

func (n DefaultNotifier) GetNotification() (Result, time.Duration, error) {
	now := n.clock.Now()
	leaderboard, err := n.repo.GetLeaderBoard(uint(now.Year()))
	if err != nil {
		return EmptyResult, 10 * time.Minute, err
	}

	msg := fmt.Sprintf(NEW_PUZZLE_FMT, AOC_URL, now.Year(), now.Day())
	result := Result{
		LeaderBoard: leaderboard,
		Messages:    []string{msg},
	}

	return result, n.GetSleep(), nil
}

func (n DefaultNotifier) GetSleep() time.Duration {
	location := time.FixedZone("", AOC_LOCATION_UTC_OFFSET*60*60)
	now := n.clock.Now().In(location)
	firstOfDecember := time.Date(now.Year(), time.December, 1, 0, 0, 0, 0, location)
	sleep := firstOfDecember.Sub(now)

	if sleep > 0 {
		return sleep
	}

	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, location)
	sleep = tomorrow.Sub(now)
	return sleep
}
