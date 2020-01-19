package main

import (
	"fmt"
	"os"

	"github.com/bowd/quip-exporter/commands"
)

var (
	buildVersion string
)

func main() {
	commands.RootCmd.Version = buildVersion

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
