package matrix

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/erdnaxeli/rudolphe/bot"
	"github.com/erdnaxeli/rudolphe/leaderboard"
	"golang.org/x/exp/slog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

const MSG_MAX_SIZE = 10

type Client struct {
	startTime time.Time

	bot    bot.Bot
	client *mautrix.Client
	roomID id.RoomID
}

type Config struct {
	HomeserverURL string
	UserID        string
	AccessToken   string
	RoomID        string
	Bot           bot.Bot
}

func NewClient(config Config) (Client, error) {
	matrixClient, err := mautrix.NewClient(
		config.HomeserverURL, id.UserID(config.UserID), config.AccessToken,
	)
	if err != nil {
		//nolint:errorlint
		return Client{}, fmt.Errorf("Error while creating client: %v", err)
	}

	client := Client{
		bot:       config.Bot,
		startTime: time.Now(),
		client:    matrixClient,
		roomID:    id.RoomID(config.RoomID),
	}
	syncer := matrixClient.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, client.onEventMessage)

	return client, nil
}

func (c Client) onEventMessage(source mautrix.EventSource, evt *event.Event) {
	evtTime := time.Unix(evt.Timestamp/1000, 0)
	if evtTime.Before(c.startTime) {
		return
	}

	content := evt.Content.AsMessage()
	result, err := c.bot.ParseCommand(content.Body)
	if err != nil {
		if errors.Is(err, bot.ErrUnknownCommand) {
			return
		}

		slog.Error("Error while parsing command", "error", err)
		return
	}

	if !result.LeaderBoard.IsEmpty() {
		err := c.sendLeaderBoard(evt.RoomID, result.LeaderBoard)
		if err != nil {
			slog.Error("Error while sending leaderboard message", "error", err)
			return
		}
	}

	for _, msg := range result.Messages {
		err := c.sendText(evt.RoomID, msg)
		if err != nil {
			slog.Error("Error while sending message", "error", err)
			return
		}
	}
}

func (b Client) sendText(roomID id.RoomID, message string) error {
	err := b.sendMessage(roomID, message, "")
	return err

}

func (b Client) sendMessage(roomID id.RoomID, message string, formattedMsg string) error {
	if formattedMsg == "" {
		formattedMsg = message
	}

	_, err := b.client.SendMessageEvent(
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

func (b Client) sendLeaderBoard(roomID id.RoomID, leaderBoard leaderboard.LeaderBoard) error {
	msg := leaderBoard.String()
	lines := strings.Split(msg, "\n")
	chunk := strings.Join(lines[0:min(len(lines), MSG_MAX_SIZE)], "\n")
	formatted := fmt.Sprintf("<pre>%s</pre>", chunk)
	err := b.sendMessage(b.roomID, chunk, formatted)
	if err != nil {
		return err
	}

	for len(lines) > MSG_MAX_SIZE {
		lines = lines[MSG_MAX_SIZE:]
		chunk := strings.Join(lines[0:min(len(lines), MSG_MAX_SIZE)], "\n")
		formatted := fmt.Sprintf("<pre>%s</pre>", chunk)
		err = b.sendMessage(b.roomID, chunk, formatted)
		if err != nil {
			return err
		}
	}

	return err
}
