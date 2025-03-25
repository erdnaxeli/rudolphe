package aoc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/erdnaxeli/rudolphe/leaderboard"
)

func (a JsonClient) GetLeaderBoard(ctx context.Context, year uint) (leaderboard.LeaderBoard, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf(
			"%s/%d/leaderboard/private/view/%s.json",
			URL, year, a.leaderBoardID,
		),
		nil,
	)
	if err != nil {
		return leaderboard.LeaderBoard{}, fmt.Errorf("error while crafting HTTP request: %v", err)
	}
	resp, err := a.client.Do(request)
	if err != nil {
		return leaderboard.LeaderBoard{}, fmt.Errorf("error during HTTP request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return leaderboard.LeaderBoard{}, ErrLeaderBoardNotFound
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return leaderboard.LeaderBoard{}, fmt.Errorf("error while reading body: %v", err)
	}

	var jsonLB jsonLeaderBoard
	err = json.Unmarshal(content, &jsonLB)
	if err != nil {
		return leaderboard.LeaderBoard{}, fmt.Errorf("error while reading json: %v", err)
	}

	lb := leaderboard.NewLeaderBoard()
	for memberID, member := range jsonLB.Members {
		user := leaderboard.NewUser()

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
				err := fmt.Errorf("invalid number of parts completed: %d", len(parts))
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
