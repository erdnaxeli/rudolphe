package aoc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/erdnaxeli/rudolphe/leaderboard"
)

func (a JsonClient) GetLeaderBoard(year uint) (leaderboard.LeaderBoard, error) {
	resp, err := a.client.Get(
		fmt.Sprintf(
			"%s/%d/leaderboard/private/view/%s.json",
			URL, year, a.leaderBoardID,
		),
	)
	if err != nil {
		return leaderboard.LeaderBoard{}, fmt.Errorf("Error during HTTP request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return leaderboard.LeaderBoard{}, ErrLeaderBoardNotFound
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return leaderboard.LeaderBoard{}, fmt.Errorf("Error while reading body: %v", err)
	}

	var jsonLB jsonLeaderBoard
	err = json.Unmarshal(content, &jsonLB)
	if err != nil {
		return leaderboard.LeaderBoard{}, fmt.Errorf("Error while reading json: %v", err)
	}

	var lb = leaderboard.NewLeaderBoard()
	for memberID, member := range jsonLB.Members {
		var user = leaderboard.NewUser()

		user.ID = memberID
		user.Name = member.Name
		user.Score = uint(member.Score)

		lb.Users[memberID] = user

		for dayNumber, parts := range member.Days {
			var part leaderboard.Part
			switch len(parts) {
			case 0:
				part = leaderboard.None
			case 1:
				part = leaderboard.Part1
			case 2:
				part = leaderboard.Part2
			default:
				err := fmt.Errorf("Invalid number of parts completed: %d", len(parts))
				return leaderboard.LeaderBoard{}, err
			}

			n, err := strconv.Atoi(dayNumber)
			if err != nil {
				// TODO
				return leaderboard.LeaderBoard{}, err
			}

			user.Days[uint8(n)] = leaderboard.Day{Done: part}
		}
	}

	return lb, nil
}
