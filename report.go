package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	//	"sync"
)

type ReportUsersCommand struct {
	Purge bool
	Ui    cli.Ui
}

func reportUsersCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &ReportUsersCommand{
		Purge: false,
		Ui:    ui,
	}, nil
}

func (c *ReportUsersCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("report", flag.ContinueOnError)
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
		fmt.Printf("error getting ldap config: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("scanning for active ldap users\n")
	lusers, err := runner.ScanForActiveUsers()
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

	for _, u := range lusers {
		for _, v := range susers {
			if v.Name == "slackbot" {
				continue
			}
			if strings.EqualFold(u.Email, v.Profile.Email) {
				// if u.Email == v.Profile.Email {
				fmt.Printf("slacker: %s,%s,%s\n", u.CN, u.Division, u.Department)
			}
		}
	}

	return 0
}

func (c *ReportUsersCommand) Help() string {
	helpText := `usage: slackers active

List all active slack users

	`
	return strings.TrimSpace(helpText)
}

func (c *ReportUsersCommand) Synopsis() string {
	return "list all live slackers"
}
