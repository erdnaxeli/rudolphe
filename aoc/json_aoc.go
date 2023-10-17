package aoc

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

type JsonClient struct {
	session       string
	leaderBoardID string

	client http.Client
}

type jsonLeaderBoard struct {
	Members map[string]member `json:"members"`
}

type member struct {
	Days  map[string]map[string]any `json:"completion_day_level"`
	Score int                       `json:"local_score"`
	Name  string                    `json:"name"`
}

func NewJsonClient(session string, leaderBoardID string) (JsonClient, error) {
	sessionCookie := http.Cookie{
		Name:     "session",
		Value:    session,
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
	}
	url, err := url.Parse(URL)
	if err != nil {
		// TODO
		return JsonClient{}, err
	}
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		// TODO
		return JsonClient{}, err
	}
	cookieJar.SetCookies(url, []*http.Cookie{&sessionCookie})

	return JsonClient{
		session:       session,
		leaderBoardID: leaderBoardID,
		client: http.Client{
			Jar:     cookieJar,
			Timeout: 30 * time.Second,
		},
	}, nil
}
