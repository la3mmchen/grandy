package app

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/la3mmchen/grandy/internal/types"
	"github.com/urfave/cli/v2"
)

/*
	Stats is the cli subcommand that collect some metrics about the input logfile.
		Prints results afterwards.

*/
func Stats(cfg *types.Config) *cli.Command {
	cmd := cli.Command{
		Name:  "stats",
		Usage: "print stats about the log events in the input.",
	}

	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "urlpath",
			Value:       "/token",
			Destination: &cfg.SearchPath,
			DefaultText: "print statistics for this path",
			Usage:       "select a search path (.e.g / or /auth) ",
		},
	}

	cmd.Action = func(c *cli.Context) error {

		err := preload(cfg)

		if err != nil {
			return err
		}

		// some maps to count
		pathOccurences := make(map[string]int)
		httpMethodOccurences := make(map[string]int)
		timeOccurences := make(map[int]int)

		fmt.Printf("`- doing stats on field [%v] \n", cfg.FileHeader[cfg.FieldInHeaderMap])

		// iterate through file
		f, err := os.Open(cfg.FileToScan)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Scan() // skip first line as it contains the header
		count := 0
		for scanner.Scan() {
			count += 1
			output, err := loglineToJson(cfg, scanner.Text())

			if err != nil {
				continue // skip log line as it was not jsonized
			}

			/*
				statistics based on the logevent field "requests"
			*/
			r := regexp.MustCompile(`^(?P<method>GET|POST|HEAD|PUT|DELETE|OPTIONS|TRACE|PATCH)\s(?P<path>.+)\sHTTP/.*$`)
			matches := r.FindStringSubmatch(fmt.Sprintf("%v", output["request"]))
			for i, name := range r.SubexpNames() {
				if i != 0 && name != "" {
					// count paths
					if name == "path" {
						u, err := url.Parse(matches[i])
						if err != nil {
							continue // do not handle unparsable times
						}
						if _, ok := pathOccurences[u.Path]; ok {
							pathOccurences[u.Path] += 1
						} else {
							pathOccurences[u.Path] = 1
						}
					} else if name == "method" {
						if _, ok := httpMethodOccurences[matches[i]]; ok {
							httpMethodOccurences[matches[i]] += 1
						} else {
							httpMethodOccurences[matches[i]] = 1
						}
					}
				}
			}

			/*
				get times of event by urlpath
			*/
			if strings.Contains(fmt.Sprintf("%v", output["request"]), cfg.SearchPath) {
				timeFromLogEvent, e := time.Parse(time.RFC3339, fmt.Sprintf("%v", output["timestamp"]))
				if e != nil {
					continue // do not handle unparsable times
				}
				if _, ok := timeOccurences[timeFromLogEvent.Hour()]; ok {
					timeOccurences[timeFromLogEvent.Hour()] += 1
				} else {
					timeOccurences[timeFromLogEvent.Hour()] = 1
				}
			}

			if count >= cfg.LineLimit && cfg.LineLimit != 0 {
				fmt.Println("`- reached LineLimit.")
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		/*
			print statistics based on the maps we created above.
		*/

		fmt.Println()
		fmt.Println("** Statistics **")

		fmt.Println("\n called base paths without params")

		for v, i := range pathOccurences {
			if strings.Contains(v, cfg.SearchPath) {
				fmt.Printf("  %s : %d \n", v, i)
			}
		}

		fmt.Println("\n used http methods")
		for v, i := range httpMethodOccurences {
			fmt.Printf("  %s : %d \n", v, i)
		}


		if (len(timeOccurences) > 0) {
			fmt.Printf("\n number of requests to [%v] by hour of day.\n", cfg.SearchPath)

			hours := make([]int, len(timeOccurences))
			i := 0
			sum := 0
			for k := range timeOccurences {
				hours[i] = k
				sum += timeOccurences[k]
				i++
			}
			sort.Ints(hours)

			for _, k := range hours {
				stars := strings.Repeat("*", (timeOccurences[k] * 100 / sum))
				fmt.Printf("	%v: %v (%v) \n", k, stars, timeOccurences[k])
			}
		} else {
			fmt.Println("\n no events found.")
		}

		return nil
	}
	return &cmd
}
