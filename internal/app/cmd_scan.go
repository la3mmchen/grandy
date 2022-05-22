package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/la3mmchen/grandy/internal/types"
	"github.com/urfave/cli/v2"
)

func Scan(cfg *types.Config) *cli.Command {
	cmd := cli.Command{
		Name:  "scan",
		Usage: "scan the configured file.",
	}
	cmd.Action = func(c *cli.Context) error {

		err := preload(cfg)

		if err != nil {
			return err
		}

		fmt.Printf("`- printing content for field [%v] \n", cfg.FileHeader[cfg.FieldInHeaderMap])

		// iterate through file
		f, err := os.Open(cfg.FileToScan)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Scan() // skip first line as we expect it to contain the header
		count := 0
		for scanner.Scan() {
			count += 1
			jsonized, err := loglineToJson(cfg, scanner.Text())
			if err != nil {
				continue // skip log line as it was not jsonized
			}
			output, err := json.MarshalIndent(jsonized, "", "  ")
			if err != nil {
				continue
			}

			fmt.Println(string(output))

			if count >= cfg.LineLimit && cfg.LineLimit != 0 {
				fmt.Println("\nReached LineLimit. Ending here.")
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		return nil
	}
	return &cmd
}
