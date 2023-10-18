package bot

import (
	"strconv"
	"time"
)

func (b bot) ParseCommand(msg string) (Result, error) {
	if msg == "!lb" {
		now := time.Now()
		year := now.Year()

		if now.Month() < 12 {
			year--
		}

		lb, err := b.leaderBoardRepo.GetLeaderBoard(uint(year))
		if err != nil {
			return EmptyResult, err
		}

		return Result{LeaderBoard: lb}, nil
	} else if len(msg) == 8 && msg[:4] == "!lb " {
		year, err := strconv.Atoi(msg[4:8])
		if err != nil {
			return Result{Messages: []string{"Année invalide"}}, nil
		}

		lb, err := b.leaderBoardRepo.GetLeaderBoard(uint(year))
		if err != nil {
			return EmptyResult, err
		}

		return Result{LeaderBoard: lb}, nil
	}

	return EmptyResult, ErrUnknownCommand
}
