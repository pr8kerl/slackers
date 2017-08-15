package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

type ScanUsersCommand struct {
	AccountId string
	Region    string
	Ui        cli.Ui
}

func scanUsersCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &ScanUsersCommand{
		AccountId: "",
		Region:    "",
		Ui:        ui,
	}, nil
}

func (c *ScanUsersCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("scan", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	ldapcfg, err := config.GetSection("ldap")
	if err != nil {
		fmt.Printf("error loading ldap config: %s\n", err)
		os.Exit(1)
	}
	slackcfg, err := config.GetSection("slack")
	if err != nil {
		fmt.Printf("error loading slack config: %s\n", err)
		os.Exit(1)
	}

	runner, err := NewLdapRunner(ldapcfg)
	if err != nil {
		fmt.Printf("error setting ldap config: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("scanning for disabled ldap users\n")
	lusers, err := runner.Scan()
	if err != nil {
		fmt.Printf("error scanning for ldap users: %s\n", err)
		os.Exit(1)
	}
	//	for _, u := range lusers {
	//		fmt.Printf("ldap: %s,%s\n", u.CN, u.Email)
	//	}

	srunner, err := NewSlackRunner(slackcfg)
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
	//	for _, v := range susers {
	//		fmt.Printf("slack: %s,%s,%t\n", v.Name, v.Profile.Email, v.Deleted)
	//	}

	for _, u := range lusers {
		for _, v := range susers {
			if u.Email == v.Profile.Email {
				if v.Name == "slackbot" {
					continue
				}
				fmt.Printf("warning: found a live slacker that should be dead: %s,%s,%t,%s,%s\n", v.Name, v.Profile.Email, v.Deleted, u.CN, u.Email)
			}
		}
	}

	return 0
}

func (c *ScanUsersCommand) Help() string {
	helpText := `usage: lusers scan

List all disabled users within the organization

	`
	return strings.TrimSpace(helpText)
}

func (c *ScanUsersCommand) Synopsis() string {
	return "scan for all disabled users within the organization"
}
