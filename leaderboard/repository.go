package leaderboard

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"
)

type Repository interface {
	GetLeaderBoard(ctx context.Context, year uint) (LeaderBoard, error)
	SaveLeaderBoard(ctx context.Context, year uint, leaderboard LeaderBoard) error
}

type SqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(filename string) (SqliteRepository, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", filename))
	if err != nil {
		slog.Error("Unable to open database", "error", err)
		return SqliteRepository{}, ErrDBOpen
	}

	return SqliteRepository{db: db}, nil
}
