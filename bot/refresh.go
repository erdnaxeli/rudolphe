package bot

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/erdnaxeli/rudolphe/leaderboard"
)

func (b bot) Refresh() (Result, time.Duration, error) {
	now := time.Now()
	var sleep time.Duration
	if time.January <= now.Month() && now.Month() <= time.November {
		sleep = 5 * time.Hour
	} else {
		sleep = MIN_REFRESH_SLEEP
	}

	remainingTimeBeforeRefresh := sleep - time.Since(b.lastRefresh)
	if remainingTimeBeforeRefresh > 0 {
		slog.Info("Refreshing too early, skipping")
		return EmptyResult, remainingTimeBeforeRefresh, nil
	}

	result := Result{}
	b.lastRefresh = now

	for year := now.Year(); year >= 2015; year-- {
		slog.Info("Updating leaderboard", "year", year)

		lb, err := b.aocClient.GetLeaderBoard(uint(year))
		if err != nil {
			slog.Error(
				"Unable to refresh leaderboard: %v",
				"year", year,
				"error", err,
			)
			break
		}

		prevLb, err := b.leaderBoardRepo.GetLeaderBoard(uint(year))
		if err != nil {
			slog.Error(
				"Unable to get previous leaderboard",
				"year", year,
				"error", err,
			)
			break
		}

		diff := lb.Difference(prevLb)

		err = b.leaderBoardRepo.SaveLeaderBoard(uint(year), lb)
		if err != nil {
			slog.Error(
				"Unable to save leaderboards: %v",
				"year", year,
				"error", err,
			)
			break
		}

		messages := b.printDiff(year, diff)
		result.Messages = append(result.Messages, messages...)
	}

	return result, sleep, nil
}

func (b bot) printDiff(year int, diff leaderboard.Diff) []string {
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
