package bot

import (
	"context"
	"time"
)

type LimitedRefresher interface {
	Refresh(ctx context.Context) (Result, time.Duration, error)
}

const (
	MIN_REFRESH_SLEEP  = 15 * time.Minute
	IDLE_REFRESH_SLEEP = 6 * time.Hour
)

// MinLimitedRefresher apply the minimum refresh limit.
//
// The minimum limit is defined by MIN_REFRESH_SLEEP.
type MinLimitedRefresher struct {
	clock     Clock
	refresher Refresher
}

func NewMinLimitedRefresher(clock Clock, refresher Refresher) MinLimitedRefresher {
	return MinLimitedRefresher{
		clock:     clock,
		refresher: refresher,
	}
}

func (m MinLimitedRefresher) Refresh(ctx context.Context) (Result, time.Duration, error) {
	now := m.clock.Now()
	return tryRefresh(ctx, now, MIN_REFRESH_SLEEP, m.refresher)
}

// MonthLimiterRefresher apply a limit depending on the current month.
//
// If we are on december, it applies the minimum limit like MinLimitedRefresher
// does. Else it applies the IDLE_REFRESH_SLEEP limit.
type MonthLimiterRefresher struct {
	clock     Clock
	refresher Refresher
}

func NewMonthLimiterRefresher(clock Clock, refresher Refresher) MonthLimiterRefresher {
	return MonthLimiterRefresher{
		clock:     clock,
		refresher: refresher,
	}
}

func (m MonthLimiterRefresher) Refresh(ctx context.Context) (Result, time.Duration, error) {
	now := m.clock.Now()
	var sleep time.Duration
	if now.Month() == time.December {
		sleep = MIN_REFRESH_SLEEP
	} else {
		sleep = IDLE_REFRESH_SLEEP
	}

	return tryRefresh(ctx, now, sleep, m.refresher)
}

func tryRefresh(ctx context.Context, now time.Time, sleep time.Duration, refresher Refresher) (Result, time.Duration, error) {
	remainingSleep := sleep - now.Sub(refresher.GetLastRefresh())
	if remainingSleep > 0 {
		return EmptyResult, remainingSleep, nil
	}

	result, err := refresher.Refresh(ctx)
	return result, sleep, err
}
