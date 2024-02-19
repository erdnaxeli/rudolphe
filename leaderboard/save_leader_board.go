package leaderboard

import "context"

func (r SqliteRepository) SaveLeaderBoard(ctx context.Context, year uint, leaderboard LeaderBoard) error {
	tx, err := r.db.Begin()
	if err != nil {
		// TODO
		return err
	}

	for _, user := range leaderboard.Users {
		_, err := tx.ExecContext(
			ctx,
			`
				INSERT OR REPLACE INTO
				users (id, name)
				VALUES
					(?, ?)
			`,
			user.ID,
			user.Name,
		)
		if err != nil {
			// TODO
			return err
		}

		row := tx.QueryRow(
			`
				INSERT OR REPLACE INTO
				scores (user_id, year, score)
				VALUES
					(?, ?, ?)
				RETURNING id
			`,
			user.ID,
			year,
			user.Score,
		)
		var scoreID int
		err = row.Scan(&scoreID)
		if err != nil {
			// TODO
			return err
		}

		for dayNumber, day := range user.Days {
			var parts uint8
			switch day.Done {
			case None:
				parts = 0
			case Part1:
				parts = 1
			case Part2:
				parts = 2
			}

			_, err = tx.Exec(
				`
					INSERT OR REPLACE INTO
					days (score_id, day, parts)
					VALUES
						(?, ?, ?)
				`,
				scoreID,
				dayNumber,
				parts,
			)
			if err != nil {
				// TODO
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		// TODO
		return err
	}

	return nil
}
