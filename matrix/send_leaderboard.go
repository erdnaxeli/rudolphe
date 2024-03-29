package matrix

import (
	"context"
	"fmt"
	"strings"

	"github.com/erdnaxeli/rudolphe/leaderboard"
	"maunium.net/go/mautrix/id"
)

func (c Client) SendLeaderBoard(ctx context.Context, leaderBoard leaderboard.LeaderBoard) error {
	err := c.sendLeaderBoard(ctx, id.RoomID(c.config.RoomID), leaderBoard)
	if err != nil {
		return ErrSendMessage
	}

	return nil
}

func (c Client) sendLeaderBoard(ctx context.Context, roomID id.RoomID, leaderBoard leaderboard.LeaderBoard) error {
	msg := leaderBoard.String()
	lines := strings.Split(msg, "\n")
	chunk := strings.Join(lines[0:min(len(lines), MSG_MAX_SIZE)], "\n")
	formatted := fmt.Sprintf("<pre>%s</pre>", chunk)
	err := c.sendMessage(ctx, roomID, chunk, formatted)
	if err != nil {
		return err
	}

	for len(lines) > MSG_MAX_SIZE {
		lines = lines[MSG_MAX_SIZE:]
		chunk := strings.Join(lines[0:min(len(lines), MSG_MAX_SIZE)], "\n")
		formatted := fmt.Sprintf("<pre>%s</pre>", chunk)
		err = c.sendMessage(ctx, roomID, chunk, formatted)
		if err != nil {
			return err
		}
	}

	return err
}
