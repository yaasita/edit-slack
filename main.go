package main

import (
	"./slack"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

const version = "0.8.1"

var api_token *string = flag.String("token", "", "api token")
var cachefile *string = flag.String("cache", "", "cache file")
var outfile *string = flag.String("outfile", "", "out file")

func main() {
	flag.Parse()
	// api token
	if *api_token == "" || *cachefile == "" || *outfile == "" {
		usage()
	}
	slk := slack.New(*api_token, *cachefile)
	// main
	if len(flag.Args()) == 0 {
		usage()
	}
	if r := regexp.MustCompile(`(history|post|search)$`); r.MatchString(flag.Args()[0]) && len(flag.Args()) < 2 {
		usage()
	}
	var stdin []byte
	if r := regexp.MustCompile(`post$`); r.MatchString(flag.Args()[0]) {
		stdin, _ = ioutil.ReadAll(os.Stdin)
	}
	var result string
	switch flag.Args()[0] {
	case "channels.list":
		result = slk.ChList()
	case "channels.history":
		result = slk.ChHistory(flag.Args()[1])
	case "channels.post":
		slk.ChannelPost(flag.Args()[1], string(stdin))
	case "users.list":
		result = slk.UserList()
	case "users.history":
		result = slk.UserHistory(flag.Args()[1])
	case "users.post":
		slk.UserPost(flag.Args()[1], string(stdin))
	case "groups.list":
		result = slk.PgList()
	case "groups.history":
		result = slk.PgHistory(flag.Args()[1])
	case "groups.post":
		slk.GroupPost(flag.Args()[1], string(stdin))
	case "search":
		result = slk.Search(flag.Args()[1])
	default:
		usage()
	}
	ioutil.WriteFile(*outfile, []byte(result+"\n"), os.ModePerm)
	slk.SaveCache()
}
func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  %s command\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `Commands:
  channels.list: Display a list of channels
  channels.history CHANNEL: Display the history of the channel
  channels.post CHANNEL: Post standard input to the channel
  users.list: Display a list of users
  users.history USER: Display history with user
  users.post USER: Send standard input to the user
  groups.list: Display a list of private groups
  groups.history GROUP: Display history of private group
  groups.post GROUP: Post standard input to private group
  search WORD: search word
`)
	fmt.Fprintf(os.Stderr, "version:\n  %s\n", version)
	os.Exit(1)
}
