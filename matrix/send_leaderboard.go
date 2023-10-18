package matrix

import "github.com/erdnaxeli/rudolphe/leaderboard"

func (c Client) SendLeaderBoard(leaderBoard leaderboard.LeaderBoard) error {
	err := c.sendLeaderBoard(c.roomID, leaderBoard)
	if err != nil {
		return ErrSendMessage
	}

	return nil
}
