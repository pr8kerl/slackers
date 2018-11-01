package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	c := cli.NewCLI("slackers", "0.1.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"disabled": disabledUsersCmdFactory,
		"report":   reportUsersCmdFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
