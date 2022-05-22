# grandy

a non-sophisticated logfile scanner writte in go. quite simple as it just parses logfiles of a certain type.

## build

```bash
$ make build
all modules verified
wrote binary to bin/grandy
```

## usage

```bash
$ ./bin/grandy --help
NAME:
   grandy - non-sophisticated logfile scanner.

USAGE:
   grandy [global options] command [command options] [arguments...]

VERSION:
   c6911bf-dirty

COMMANDS:
   scan     scan the configured file.
   stats    print stats about the log events in the input.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --field value        Define the field you want to see. (default: "message")
   --file value         Define the file you want to crawl. (default: "data/example.log")
   --separator value    Default separator in between the fields (default: ",")
   --limit value        Set a limit for lines to be parsed. Set to 0 to get all lines. (default: limit of lines to parse)
   --help, -h           show help (default: false)
   --print-version, -V  print only the version (default: false)
```

## example

```bash
