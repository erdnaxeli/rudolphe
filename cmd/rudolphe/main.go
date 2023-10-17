package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/erdnaxeli/rudolphe/aoc"
	"github.com/erdnaxeli/rudolphe/leaderboard"
	"github.com/erdnaxeli/rudolphe/matrix"
)

type Config struct {
	LeaderBoardID string
	SessionCookie string
	HomeServerURL string
	UserID        string
	AccessToken   string
	RoomID        string
}

func main() {
	repo, err := leaderboard.NewSqliteRepository("db.sqlite")
	if err != nil {
		log.Fatalf("Unable to create repository: %v", err)
	}

	content, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error while reading config: %v", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Error while parsing config: %v", err)
	}

	bot, err := matrix.NewBot(
		matrix.Config{
			HomeserverURL:         config.HomeServerURL,
			UserID:                config.UserID,
			AccessToken:           config.AccessToken,
			RoomID:                config.RoomID,
			LeaderBoardRepository: repo,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = startUpdater(repo, bot, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(bot.StartSync())
	//mClient.SendText()
}

func startUpdater(repo leaderboard.Repository, bot matrix.Bot, config Config) error {
	client, err := aoc.NewJsonClient(
		config.SessionCookie,
		config.LeaderBoardID,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			now := time.Now()
			for year := now.Year(); year >= 2015; year-- {
				slog.Info("Updating leaderboard", "year", year)

				lb, err := client.GetLeaderBoard(uint(year))
				if err != nil {
					slog.Error(
						"Unable to refresh leaderboard: %v",
						"year", year,
						"error", err,
					)
					break
				}

				prevLb, err := repo.GetLeaderBoard(uint(year))
				if err != nil {
					slog.Error(
						"Unable to get previous leaderboard",
						"year", year,
						"error", err,
					)
					break
				}

				diff := lb.Difference(prevLb)
				err = printDiff(year, diff, bot)
				if err != nil {
					slog.Error(
						"Unable to send diff to matrix",
						"year", year,
						"error", err,
					)
				}

				err = repo.SaveLeaderBoard(uint(year), lb)
				if err != nil {
					slog.Error(
						"Unable to save leaderboards: %v",
						"year", year,
						"error", err,
					)
					break
				}
			}

			var sleep time.Duration
			if time.January <= now.Month() && now.Month() <= time.November {
				sleep = 5 * time.Hour
			} else {
				sleep = 15 * time.Minute
			}

			slog.Info("Going to sleep before next refresh", "sleep", sleep)
			time.Sleep(sleep)
		}
	}()

	return nil
}

func printDiff(year int, diff leaderboard.Diff, bot matrix.Bot) error {
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
			err := bot.SendText(builder.String())
			if err != nil {
				return err
			}
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

		err := bot.SendText(builder.String())
		if err != nil {
			return err
		}
	}

	return nil
}
