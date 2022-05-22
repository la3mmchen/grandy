package app

import (
	"github.com/la3mmchen/grandy/internal/types"
	"github.com/urfave/cli/v2"
)

// CreateApp builds the cli app by enriching the
// urfave/cli app struct with our params, flags, and commands.
// returns a pointer to a cli.App struct
func CreateApp(cfg *types.Config) *cli.App {
	cliFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "field",
			Value:       "message",
			Usage:       "Define the field you want to see.",
			Destination: &cfg.FieldToPrint,
		},
		&cli.StringFlag{
			Name:        "file",
			Value:       "data/example.log",
			Usage:       "Define the file you want to crawl.",
			Destination: &cfg.FileToScan,
		},
		&cli.StringFlag{
			Name:        "separator",
			Value:       ",",
			Usage:       "Default separator in between the fields",
			Destination: &cfg.Separator,
		},
		&cli.IntFlag{
			Name:        "limit",
			Value:       2,
			Destination: &cfg.LineLimit,
			DefaultText: "limit of lines to parse",
			Usage:       "Set a limit for lines to be parsed. Set to 0 to get all lines.",
		},
	}

	cliFlags = append(cliFlags)

	app := cli.App{
		Name:    cfg.AppName,
		Usage:   cfg.AppUsage,
		Version: cfg.AppVersion,
		Commands: []*cli.Command{
			Scan(cfg),
			Stats(cfg),
		},
		Flags: cliFlags,
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	return &app
}
