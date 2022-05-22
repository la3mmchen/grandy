package main

import (
	"fmt"
	"os"

	"github.com/la3mmchen/grandy/internal/app"
	"github.com/la3mmchen/grandy/internal/types"
)

var (
	// AppVersion to be injected during build time
	AppVersion string
)

func main() {

	cfg := types.Config{
		AppName:          "grandy",
		AppUsage:         "non-sophisticated logfile scanner.",
		AppVersion:       AppVersion,
		FileHeader:       make(map[int]string),
		FieldInHeaderMap: -1,
	}

	app := app.CreateApp(&cfg)

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}
}
