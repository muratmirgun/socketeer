package main

import (
	"github.com/muratmirgun/socketeer/internal/commands"
	"github.com/muratmirgun/socketeer/internal/version"
)

var (
	buildVersion = "dev"
	buildCommit  = "unknown"
	buildDate    = "unknown"
)

func main() {
	// Set version information
	version.Version = buildVersion
	version.Commit = buildCommit
	version.Date = buildDate
	
	commands.Execute()
}
