package main

import (
	"github.com/muratmirgun/socketeer/internal/commands"
	"github.com/muratmirgun/socketeer/internal/version"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Set version information
	version.Version = version
	version.Commit = commit
	version.Date = date
	
	commands.Execute()
}
