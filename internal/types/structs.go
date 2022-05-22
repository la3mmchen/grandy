package types

// Config contains the structured config during our lifetime
type Config struct {
	// basic settings for the cli itself
	AppName    string
	AppUsage   string
	AppVersion string

	// content to be injected via flags
	FieldToPrint string
	FileToScan   string
	Separator    string
	LineLimit    int
	SearchPath   string

	// internal fields
	FileHeader       map[int]string
	FieldInHeaderMap int
}
