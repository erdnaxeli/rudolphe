package leaderboard

import (
	"golang.org/x/exp/slog"
)

func (r SqliteRepository) GetLeaderBoard(year uint) (LeaderBoard, error) {
	rows, err := r.db.Query(
		`
			SELECT
				users.id,
				users.name,
				scores.score,
				days.day,
				days.parts
			FROM users
			JOIN scores ON
				users.id = scores.user_id
				AND scores.year = ?
			JOIN days ON
				scores.id = days.score_id
		`,
		year,
	)
	if err != nil {
		// TODO
		return emptyLeaderBoard, err
	}

	var (
		userID   string
		userName string
		score    int
		day      uint8
		parts    uint8
	)
	leaderboard := NewLeaderBoard()
	for rows.Next() {
		err := rows.Scan(&userID, &userName, &score, &day, &parts)
		if err != nil {
			return emptyLeaderBoard, err
		}

		user, ok := leaderboard.Users[userID]
		if !ok {
			// adding new user
			user = NewUser()
			user.ID = userID
			user.Name = userName
			user.Score = uint(score)
			leaderboard.Users[userID] = user
		}

		var part Part
		if parts == 0 {
			part = None
		} else if parts == 1 {
			part = Part1
		} else if parts == 2 {
			part = Part2
		} else {
			slog.Warn(
				"Invalid part count for user",
				"user", userID,
				"parts", parts,
			)
			part = Part2
		}

		user.Days[day] = Day{Done: part}
	}

	return leaderboard, nil
}
