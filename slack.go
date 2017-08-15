package main

import (
	"github.com/nlopes/slack"
	"gopkg.in/ini.v1"
)

type SlackRunner struct {
	token string
}

func (s *SlackRunner) Scan() (map[string]slack.User, error) {
	api := slack.New(s.token)
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)
	users, err := api.GetUsers()
	if err != nil {
		return nil, err
	}
	active := make(map[string]slack.User)
	for _, u := range users {
		if !u.Deleted {
			active[u.Profile.Email] = u
		}
	}
	return active, nil
}

func NewSlackRunner(cfg *ini.Section) (*SlackRunner, error) {
	token, err := cfg.GetKey("api_token")
	if err != nil {
		return nil, err
	}
	return &SlackRunner{
		token: token.Value(),
	}, nil
}
