package leaderboard

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/exp/slices"
)

type Part int

const (
	None Part = iota
	Part1
	Part2
)

type LeaderBoard struct {
	Users map[string]User
}

type User struct {
	Days  map[uint8]Day
	ID    string
	Name  string
	Score uint
}

func (u User) GetName() string {
	if u.Name != "" {
		return u.Name
	}

	return fmt.Sprintf("anon %s", u.ID)
}

func (u User) GetNameWithoutHL() string {
	// It joins with a "zero width space" char.
	return strings.Join(strings.Split(u.GetName(), ""), "\u200B")
}

type Day struct {
	Done Part
}

var emptyLeaderBoard = LeaderBoard{}

func NewLeaderBoard() LeaderBoard {
	return LeaderBoard{
		Users: make(map[string]User),
	}
}

func NewUser() User {
	return User{
		Days: make(map[uint8]Day),
	}
}

func (l LeaderBoard) IsEmpty() bool {
	return len(l.Users) == 0
}

func (l LeaderBoard) String() string {
	var users []User
	for _, user := range l.Users {
		if user.Score > 0 {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		return "No one has done anything yet."
	}

	slices.SortFunc(users, func(a User, b User) int {
		if a.Score < b.Score {
			return 1
		} else if a.Score == b.Score {
			return 0
		} else {
			return -1
		}
	})

	// the size of the max position count in the leaderboard
	// 15 users: size of 2 chars
	maxPositionSize := len(fmt.Sprint(len(users)))

	maxScore := users[0].Score
	// the size of the max score
	// 658: size of 3 chars
	maxScoreSize := len(fmt.Sprint(maxScore))

	var builder strings.Builder

	// header
	//
	// It looks like this:
	//
	//                  1111111111222222
	//         1234567890123456789012345
	//  1) 658 *+*                   +
	// ...
	// 10)  76 *

	// the spaces before the numbers tens
	printSpaces(&builder, maxPositionSize+2+maxScoreSize)
	// the tens
	fmt.Fprint(&builder, "          1111111111222222\n")

	// the spaces before the numbers units
	printSpaces(&builder, maxPositionSize+2+maxScoreSize)
	// the units
	fmt.Fprint(&builder, " 1234567890123456789012345\n")

	position := 1
	for _, user := range users {
		printSpaces(&builder, maxPositionSize-len(fmt.Sprint(position)))
		fmt.Fprintf(&builder, "%d) ", position)
		printSpaces(&builder, maxScoreSize-len(fmt.Sprint(user.Score)))
		fmt.Fprintf(&builder, "%d ", user.Score)

		for i := 1; i <= 25; i++ {
			day, ok := user.Days[uint8(i)]
			if !ok {
				day.Done = None
			}

			switch day.Done {
			case None:
				fmt.Fprint(&builder, " ")
			case Part1:
				fmt.Fprint(&builder, "·")
			case Part2:
				fmt.Fprint(&builder, "×")
			}
		}

		fmt.Fprint(&builder, " ")
		fmt.Fprintln(&builder, user.GetNameWithoutHL())
		position++
	}

	return builder.String()
}

func printSpaces(to io.Writer, count int) {
	for i := 0; i < count; i++ {
		fmt.Fprint(to, " ")
	}
}
