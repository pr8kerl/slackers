package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/nlopes/slack"
)

type DisabledUsersCommand struct {
	Purge        bool
	OutputFormat string
	Ui           cli.Ui
}

func disabledUsersCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &DisabledUsersCommand{
		Purge:        false,
		OutputFormat: "stdout",
		Ui:           ui,
	}, nil
}

func (c *DisabledUsersCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("disabled", flag.ContinueOnError)
	cmdFlags.BoolVar(&c.Purge, "purge", false, "deactivate live slackers that are disabled in AD")
	cmdFlags.StringVar(&c.OutputFormat, "o", "stdout", "desired output format")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	ldapcfg, err := config.GetSection("ldap")
	if err != nil {
		fmt.Printf("error loading ldap config: %s\n", err)
		os.Exit(1)
	}

	runner, err := NewLdapRunner(ldapcfg)
	if err != nil {
		fmt.Printf("error setting ldap config: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("scanning for disabled ldap users\n")
	lusers, err := runner.ScanForDisabledUsers()
	if err != nil {
		fmt.Printf("error scanning for ldap users: %s\n", err)
		os.Exit(1)
	}

	srunner, err := NewSlackRunner()
	if err != nil {
		fmt.Printf("error setting slack config: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("scanning for active slack users\n")
	susers, err := srunner.Scan()
	if err != nil {
		fmt.Printf("error scanning for slack users: %s\n", err)
		os.Exit(1)
	}
	// stack for targeted users
	var tusers []slack.User

	for _, u := range lusers {
		for _, v := range susers {
			if v.Name == "slackbot" {
				continue
			}
			if v.IsBot == true {
				continue
			}
			if strings.EqualFold(u.Email, v.Profile.Email) {
				if c.OutputFormat == "stdout" {
					fmt.Printf("warning: found a live slacker that should be dead: %s,%s,%s\n", v.Name, v.Profile.Email, v.ID)
				}
				if c.OutputFormat == "json" {
					tusers = append(tusers, v)
				}
				if c.Purge {
					resp, err := srunner.DeleteUser(&v)
					if err != nil {
						fmt.Printf("error deleting slacker: %s\n", err)
					}
					fmt.Printf("killed slacker: %s %s\n", v.Name, resp.Status)
				}
			}
		}
	}
	if c.OutputFormat == "json" {
		js, _ := json.Marshal(tusers)
		fmt.Println(string(js))
	}
	return 0
}

func (c *DisabledUsersCommand) Help() string {
	helpText := `usage: slackers disabled

List all disabled AD users that are live slackers

Options:
            -purge          deactivate live slackers that are deactived in AD


	`
	return strings.TrimSpace(helpText)
}

func (c *DisabledUsersCommand) Synopsis() string {
	return "scan for all disabled AD users that are live slackers"
}
