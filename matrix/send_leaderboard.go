package matrix

import "github.com/erdnaxeli/rudolphe/leaderboard"

func (b Bot) SendLeaderBoard(leaderBoard leaderboard.LeaderBoard) error {
	err := b.sendLeaderBoard(b.roomID, leaderBoard)
	if err != nil {
		return ErrSendMessage
	}

	return nil
}
