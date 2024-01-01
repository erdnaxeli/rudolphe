package bot

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/erdnaxeli/rudolphe/aoc"
	"github.com/erdnaxeli/rudolphe/leaderboard"
)

type Refresher interface {
	// Refresh all leaderboards.
	//
	// If any update happened, it returns the corresponding messages.
	Refresh() (Result, error)

	// Return the last refresh time.
	GetLastRefresh() time.Time
}

type DefaultRefresher struct {
	aocClient aoc.Client
	clock     Clock
	repo      leaderboard.Repository

	lastRefresh time.Time
}

func NewRefresher(
	aocClient aoc.Client,
	clock Clock,
	repo leaderboard.Repository,
) *DefaultRefresher {
	return &DefaultRefresher{
		aocClient: aocClient,
		clock:     clock,
		repo:      repo,
	}
}

func (r *DefaultRefresher) GetLastRefresh() time.Time {
	return r.lastRefresh
}

func (r *DefaultRefresher) Refresh() (Result, error) {
	result := Result{}
	now := time.Now()
	r.lastRefresh = now

	maxYear := now.Year()
	if now.Month() < time.December {
		maxYear--
	}

	for year := maxYear; year >= 2015; year-- {
		slog.Info("Updating leaderboard", "year", year)

		lb, err := r.aocClient.GetLeaderBoard(uint(year))
		if err != nil {
			slog.Error(
				"Unable to refresh leaderboard: %v",
				"year", year,
				"error", err,
			)
			break
		}

		prevLb, err := r.repo.GetLeaderBoard(uint(year))
		if err != nil {
			slog.Error(
				"Unable to get previous leaderboard",
				"year", year,
				"error", err,
			)
			break
		}

		diff := lb.Difference(prevLb)

		err = r.repo.SaveLeaderBoard(uint(year), lb)
		if err != nil {
			slog.Error(
				"Unable to save leaderboards: %v",
				"year", year,
				"error", err,
			)
			break
		}

		messages := r.printDiff(year, diff)
		result.Messages = append(result.Messages, messages...)
	}

	return result, nil
}

func (r *DefaultRefresher) printDiff(year int, diff leaderboard.Diff) []string {
	var messages []string

	for _, user := range diff.Users {
		var builder strings.Builder

		currentYear := time.Now().Year()
		if time.Now().Month() < time.December {
			currentYear--
		}
		if year != currentYear {
			fmt.Fprintf(&builder, "[%d] ", year)
		}

		if user.IsNew {
			fmt.Fprintf(
				&builder,
				"Un nouveau concurrent entre dans la place avec %d points, bienvenue à %s !",
				user.NewScore,
				user.GetNameWithoutHL(),
			)
			messages = append(messages, builder.String())
			continue
		}

		fmt.Fprintf(&builder, "%s vient juste de compléter ", user.GetNameWithoutHL())

		var daysMsg []string
		for dayNumber, parts := range user.Days {
			switch parts {
			case leaderboard.DiffPart1:
				daysMsg = append(
					daysMsg,
					fmt.Sprintf("la partie 1 du jour %d", dayNumber),
				)
			case leaderboard.DiffPart2:
				daysMsg = append(
					daysMsg,
					fmt.Sprintf("la partie 2 du jour %d", dayNumber),
				)
			case leaderboard.DiffPartBoth:
				daysMsg = append(
					daysMsg,
					fmt.Sprintf("le jour %d", dayNumber),
				)
			}
		}

		if len(daysMsg) > 2 {
			for _, m := range daysMsg[:len(daysMsg)-2] {
				fmt.Fprint(&builder, m, ", ")
			}

			fmt.Fprint(&builder, daysMsg[len(daysMsg)-2], " et ", daysMsg[len(daysMsg)-1])
		} else if len(daysMsg) == 2 {
			fmt.Fprint(&builder, daysMsg[len(daysMsg)-2], " et ", daysMsg[len(daysMsg)-1])
		} else {
			fmt.Fprint(&builder, daysMsg[0])
		}

		fmt.Fprintf(&builder, " (+%d points)", user.NewScore)

		messages = append(messages, builder.String())
	}

	return messages
}
