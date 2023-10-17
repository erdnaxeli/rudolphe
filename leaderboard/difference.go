package leaderboard

type Diff struct {
	Users []DiffUser
}

type DiffUser struct {
	Days     map[uint8]NewPart
	ID       string
	IsNew    bool
	Name     string
	NewScore uint
}

func (d DiffUser) GetNameWithoutHL() string {
	return User{
		ID:   d.ID,
		Name: d.Name,
	}.GetNameWithoutHL()
}

type NewPart int

const (
	DiffPart1 NewPart = iota
	DiffPart2
	DiffPartBoth
)

func (d Diff) String() string {
	return ""
}

// l - p
//
// If user is in l but not p, the diff is +user.
func (currentLb LeaderBoard) Difference(previousLb LeaderBoard) Diff {
	d := Diff{}

	for userID, user := range currentLb.Users {
		previousUser, ok := previousLb.Users[userID]

		if user.Score == previousUser.Score {
			// nothing changed
			continue
		}

		dUser := DiffUser{
			Days:     make(map[uint8]NewPart),
			ID:       user.ID,
			IsNew:    !ok,
			Name:     user.Name,
			NewScore: user.Score - previousUser.Score,
		}

		for dayNumber, day := range user.Days {
			var newPart NewPart
			previousDay := previousUser.Days[dayNumber]

			if day.Done == previousDay.Done {
				// nothing changed
				continue
			}

			if day.Done == Part2 {
				if previousDay.Done == None {
					newPart = DiffPartBoth
				} else {
					newPart = DiffPart2
				}
			} else {
				newPart = DiffPart1
			}

			dUser.Days[dayNumber] = newPart
		}

		if len(dUser.Days) > 0 {
			// When a user join the leader board the score is computed for every user
			// so we can see a diff, but actually no days changed.
			d.Users = append(d.Users, dUser)
		}
	}

	return d
}
