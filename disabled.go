package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	//	"sync"
)

type DisabledUsersCommand struct {
	Purge bool
	Ui    cli.Ui
}

func disabledUsersCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &DisabledUsersCommand{
		Purge: false,
		Ui:    ui,
	}, nil
}

func (c *DisabledUsersCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("disabled", flag.ContinueOnError)
	cmdFlags.BoolVar(&c.Purge, "purge", false, "deactivate live slackers that are disabled in AD")
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
	// for _, u := range lusers {
	// 	fmt.Printf("ldap: %s,%s\n", u.CN, u.Email)
	// }

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
				fmt.Printf("warning: found a live slacker that should be dead: %s,%s,%s\n", v.Name, v.Profile.Email, v.ID)

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

	return 0
}

func (c *DisabledUsersCommand) Help() string {
	helpText := `usage: slackers disabled

List all disabled AD users that are live slackers

Options:
            -purge          deactivate live slackers that are deactive in AD


	`
	return strings.TrimSpace(helpText)
}

func (c *DisabledUsersCommand) Synopsis() string {
	return "scan for all disabled AD users that are live slackers"
}
