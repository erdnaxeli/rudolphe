package matrix

import (
	"errors"
	"fmt"
	"time"

	"github.com/erdnaxeli/rudolphe/bot"
	"golang.org/x/exp/slog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

const MSG_MAX_SIZE = 10

type Client struct {
	startTime time.Time

	client *mautrix.Client
	config Config
}

type Config struct {
	HomeserverURL string
	UserID        string
	AccessToken   string
	RoomID        string
	MessageParser bot.MessageParser
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
		startTime: time.Now(),
		config:    config,
		client:    matrixClient,
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
	result, err := c.config.MessageParser.ParseMessage(content.Body)
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
