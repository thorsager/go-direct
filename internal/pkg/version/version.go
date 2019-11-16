package version

var ( // metadata is extra build time data
	version  = "v0.1"
	metadata = "unreleased"
	// gitCommit is the git sha1
	gitCommit = ""
	// gitTreeState is the state of the git tree
	gitTreeState = ""
)

// GetVersion returns the semver string of the version
func GetVersion() string {
	if //noinspection GoBoolExpressions
	metadata == "" {
		return version
	}
	return version + "+" + metadata + "+" + gitCommit + "+" + gitTreeState
}
