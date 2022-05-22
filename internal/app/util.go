package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/la3mmchen/grandy/internal/types"
)

/*
	preload common stuff that might be useful for all commands

	return error in case somethin is broken.
*/
func preload(cfg *types.Config) error {

	// check if input file exists
	if _, err := os.Stat(cfg.FileToScan); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("- does not exists. Exiting.\n")
		return err
	}

	// get headers
	err := extractHeader(cfg)

	if err != nil {
		return err
	}
	if cfg.LineLimit != 0 {
		fmt.Printf("`- running with line limit of [%v] \n", cfg.LineLimit)
	}

	return nil

}

/*
	extractHeader treads the first line of the file as header that describes the rows.
		splits via the configured separator.
		removes leading and trailing quotes.
		finds the field the user wants to see and saves it.

		throws an exception if parsing went horribly wrong.
*/
func extractHeader(cfg *types.Config) error {
	f, err := os.Open(cfg.FileToScan)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// TODO: this is very hacky. it won't work if there are additional whitespaces
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	items := strings.Split(scanner.Text(), cfg.Separator)

	for i, v := range items {
		v = strings.Trim(v, "\"") //TODO: error-handling here please
		cfg.FileHeader[i] = v
		if v == cfg.FieldToPrint {
			cfg.FieldInHeaderMap = i
		}
	}

	if cfg.FieldInHeaderMap < 0 {
		return fmt.Errorf("field [%v] not found in header line [%v].", cfg.FieldToPrint, items)
	} else {
		fmt.Printf("`- found header fields %v \n", cfg.FileHeader)
	}

	return nil
}

/*
	loglineToJson parses a logline, extracts the wanted column and returns a map that contains the json.

	returns the map.
	if string is no valid json does nothing.

	hint: the rework of the found string into a parsable json is handcrafted to some example files.
				this is no good and robust code.
*/

func loglineToJson(cfg *types.Config, text string) (map[string]any, error) {
	tmpResult := make(map[int]string)

	// iterate over the input line and cut the same number of chunks as we found header
	// 		this must be done as in the provided input file the , is used a delimiter in between the columns and also in the json string itself
	for i := 0; i < len(cfg.FileHeader)-1; i++ {
		tmpResult[i] = strings.Trim(text[:strings.IndexByte(text, ',')], "\"")
		text = text[strings.IndexByte(text, ',')+1:]
	}
	// parse the remaing input string into the last item of the array
	tmpResult[len(cfg.FileHeader)-1] = text

	// strip double double quotes and leading/trailing quotes from the example file
	foundResult := strings.Replace(tmpResult[cfg.FieldInHeaderMap], "\"\"", "\"", -1)
	foundResult = strings.Trim(foundResult, "\"")

	// create a map that contains the json
	var jsn map[string]any
	if err := json.Unmarshal([]byte(foundResult), &jsn); err != nil {
		return jsn, err
	}
	return jsn, nil
}
