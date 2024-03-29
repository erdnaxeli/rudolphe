package matrix

import (
	"context"

	"golang.org/x/exp/slog"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func (c Client) SendText(ctx context.Context, text string) error {
	err := c.sendText(ctx, id.RoomID(c.config.RoomID), text)
	if err != nil {
		return ErrSendMessage
	}

	return nil
}

func (c Client) sendText(ctx context.Context, roomID id.RoomID, message string) error {
	err := c.sendMessage(ctx, roomID, message, "")
	return err
}

func (c Client) sendMessage(ctx context.Context, roomID id.RoomID, message string, formattedMsg string) error {
	if formattedMsg == "" {
		formattedMsg = message
	}

	_, err := c.client.SendMessageEvent(
		ctx,
		roomID,
		event.EventMessage,
		map[string]string{
			"body":           message,
			"format":         "org.matrix.custom.html",
			"formatted_body": formattedMsg,
			"msgtype":        "m.text",
		},
	)
	if err != nil {
		slog.Error("Error while sending message: %v", err)
		return err
	}

	return nil
}
