package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/nlopes/slack"
)

type SlackRunner struct {
	token string
}

type slackResponse struct {
	Status string
	Header http.Header
	Body   string
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

func (s *SlackRunner) DeleteUser(u *slack.User) (*slackResponse, error) {

	url := "https://api.slack.com/scim/v1/Users/"
	url += u.ID
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	bearer := "Bearer " + s.token
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return &slackResponse{
		Status: resp.Status,
		Header: resp.Header,
		Body:   string(body),
	}, nil
}

func NewSlackRunner() (*SlackRunner, error) {
	token := os.Getenv("SLACK_API_TOKEN")
	if token == "" {
		return nil, errors.New("missing SLACK_API_TOKEN environment variable")
	}
	return &SlackRunner{
		token: token,
	}, nil
}
