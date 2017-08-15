package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

func main() {

	c := cli.NewCLI("lusermgr", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"scan": scanUsersCmdFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
