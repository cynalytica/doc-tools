package flags

const (
	ConfigFile           = "config-file"
	ConfigDebugLevelFlag = "log-level"
	ConfigDebugLevelEnv  = "LOG_LEVEL"

	// common configs

	Location  = "location"
	Title     = "title"
	Subtitle  = "subtitle"
	Regex     = "regex"
	RegexFile = "regex-file"

	// TOC configs

	IncludeContent = "content"
	Description    = "description"

	// PDF configs

	Output     = "output"
	Version    = "version"
	CommitHash = "commit"
)
