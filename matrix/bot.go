package matrix

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/erdnaxeli/rudolphe/leaderboard"
	"golang.org/x/exp/slog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

const MSG_MAX_SIZE = 10

type Bot struct {
	startTime time.Time

	client          *mautrix.Client
	leaderBoardRepo leaderboard.Repository
	roomID          id.RoomID
}

type Config struct {
	HomeserverURL         string
	UserID                string
	AccessToken           string
	RoomID                string
	LeaderBoardRepository leaderboard.Repository
}

func NewBot(config Config) (Bot, error) {
	client, err := mautrix.NewClient(
		config.HomeserverURL, id.UserID(config.UserID), config.AccessToken,
	)
	if err != nil {
		//nolint:errorlint
		return Bot{}, fmt.Errorf("Error while creating client: %v", err)
	}

	bot := Bot{
		startTime:       time.Now(),
		client:          client,
		leaderBoardRepo: config.LeaderBoardRepository,
		roomID:          id.RoomID(config.RoomID),
	}
	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, bot.onEventMessage)

	return bot, nil
}

func (b Bot) onEventMessage(source mautrix.EventSource, evt *event.Event) {
	evtTime := time.Unix(evt.Timestamp/1000, 0)
	if evtTime.Before(b.startTime) {
		return
	}

	content := evt.Content.AsMessage()

	if content.Body == "!lb" {
		now := time.Now()
		year := now.Year()

		if now.Month() < 12 {
			year--
		}

		lb, err := b.leaderBoardRepo.GetLeaderBoard(uint(year))
		if err != nil {
			_ = b.sendText(evt.RoomID, err.Error())
			return
		}

		_ = b.sendLeaderBoard(evt.RoomID, lb)
	} else if len(content.Body) == 8 && content.Body[:4] == "!lb " {
		year, err := strconv.Atoi(content.Body[4:8])
		if err != nil {
			_ = b.sendText(evt.RoomID, "Année invalide")
			return
		}

		lb, err := b.leaderBoardRepo.GetLeaderBoard(uint(year))
		if err != nil {
			_ = b.sendText(evt.RoomID, err.Error())
			return
		}

		_ = b.sendLeaderBoard(evt.RoomID, lb)
	}
}

func (b Bot) sendText(roomID id.RoomID, message string) error {
	err := b.sendMessage(roomID, message, "")
	return err

}

func (b Bot) sendMessage(roomID id.RoomID, message string, formattedMsg string) error {
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

func (b Bot) sendLeaderBoard(roomID id.RoomID, leaderBoard leaderboard.LeaderBoard) error {
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
