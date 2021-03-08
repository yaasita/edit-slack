package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/yaasita/edit-slack/editslack"
)

const version = "1.0.1"

func main() {
	options := flag.NewFlagSet("options", flag.ExitOnError)
	cache_file := options.String("cache", "", "cache file")
	chname := options.String("chname", "", "channel name")
	file_path := options.String("file", "", "file path")
	file_url := options.String("furl", "", "file download url")
	out_file := options.String("out", "", "out file")
	token := options.String("token", "", "slack token")
	search_word := options.String("word", "", "search word")
	ts := options.String("ts", "", "time stamp")
	if len(os.Args) < 2 {
		usage(options, 1)
	}
	switch os.Args[1] {
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "help":
		usage(options, 0)
	}
	if len(os.Args) < 3 {
		usage(options, 1)
	}
	options.Parse(os.Args[3:])
	if *token == "" || *cache_file == "" {
		usage(options, 1)
	}
	slk := editslack.New(*token, *cache_file)
	switch os.Args[1] {
	case "conversations":
		switch os.Args[2] {
		case "list":
			slk.ConversationList(*out_file)
		case "history":
			slk.ConversationHistory(*chname, *out_file)
		case "post":
			stdin, err := ioutil.ReadAll(os.Stdin)
			check(err, "stdin read error")
			slk.ConversationPost(*chname, string(stdin))
		case "join":
			slk.ConversationJoin(*chname)
		case "leave":
			slk.ConversationLeave(*chname)
		default:
			usage(options, 1)
		}
	case "replies":
		switch os.Args[2] {
		case "history":
			slk.RepliesHistory(*chname, *ts, *out_file)
		case "post":
			stdin, err := ioutil.ReadAll(os.Stdin)
			check(err, "stdin read error")
			slk.RepliesPost(*chname, *ts, string(stdin))
		default:
			usage(options, 1)
		}
	case "files":
		switch os.Args[2] {
		case "download":
			slk.FileDownload(*file_url, *out_file)
		case "upload":
			slk.FileUpload(*file_path, *chname, *ts)
		default:
			usage(options, 1)
		}
	case "search":
		switch os.Args[2] {
		case "all":
			slk.SearchAll(*search_word, *out_file)
		default:
			usage(options, 1)
		}
	default:
		usage(options, 1)
	}
	slk.SaveCache()
}
func usage(options *flag.FlagSet, exitcode int) {
	fmt.Println(os.Args[0], "<subcommand> <options>")
	fmt.Println(
		`subcommand: 
  conversations list    --token TOKEN --cache CACHEFILE                         --out OUTFILE
  conversations history --token TOKEN --cache CACHEFILE --chname CHNAME         --out OUTFILE 
  conversations post    --token TOKEN --cache CACHEFILE --chname CHNAME 
  conversations join    --token TOKEN --cache CACHEFILE --chname CHNAME 
  conversations leave   --token TOKEN --cache CACHEFILE --chname CHNAME 
  replies history       --token TOKEN --cache CACHEFILE --chname CHNAME --ts TS --out OUTFILE
  replies post          --token TOKEN --cache CACHEFILE --chname CHNAME --ts TS
  files download        --token TOKEN --cache CACHEFILE --furl URL              --out OUTFILE
  files upload          --token TOKEN --cache CACHEFILE --file FILEPATH         --out OUTFILE
  files upload          --token TOKEN --cache CACHEFILE --file FILEPATH --ts TS --out OUTFILE
  search all            --token TOKEN --cache CACHEFILE --word WORD             --out OUTFILE
  version
  help
options:`)
	options.PrintDefaults()
	os.Exit(exitcode)
}
func check(e error, s string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, s)
		log.Fatal(e)
	}
}
