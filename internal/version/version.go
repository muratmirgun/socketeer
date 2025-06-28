package version

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

// GetVersion returns the current version
func GetVersion() string {
	return Version
}

// GetCommit returns the current commit hash
func GetCommit() string {
	return Commit
}

// GetDate returns the build date
func GetDate() string {
	return Date
} 