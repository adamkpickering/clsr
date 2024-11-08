package main

import "github.com/adamkpickering/clsr/cmd"

var version = "development"

func main() {
	cmd.SetVersionInfo(version)
	cmd.Execute()
}
