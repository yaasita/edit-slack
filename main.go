package main

import (
	"./slack_history"
	"./slack_list"
	"./slack_post"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var api_token *string = flag.String("token", "", "api token")
var cachefile *string = flag.String("cache", "", "cache file")
var outfile *string = flag.String("outfile", "", "out file")

func main() {
	flag.Parse()
	// api token
	if *api_token == "" {
		usage()
	}
	// cache file
	hash_file := ""
	if *cachefile == "" {
		fw, err := ioutil.TempFile("", "slack")
		if err != nil {
			panic(err)
		}
		hash_file = fw.Name()
	} else {
		hash_file = *cachefile
	}
	dict := slack_list.List(*api_token, hash_file)
	// out file
	if *outfile == "" {
		usage()
	}
	// main
	if len(flag.Args()) == 0 {
		usage()
	}
	if r := regexp.MustCompile(`(history|post)$`); r.MatchString(flag.Args()[0]) && len(flag.Args()) < 2 {
		usage()
	}
	var stdin []byte
	if r := regexp.MustCompile(`post$`); r.MatchString(flag.Args()[0]) {
		stdin, _ = ioutil.ReadAll(os.Stdin)
	}
	var b string
	switch flag.Args()[0] {
	case "channels.list":
		b = slack_history.ChList(dict)
	case "channels.history":
		b = slack_history.ChHistory(*api_token, dict, flag.Args()[1])
	case "channels.post":
		b = slack_post.ChannelPost(*api_token, dict, flag.Args()[1], string(stdin))
	case "users.list":
		b = slack_history.UserList(dict)
	case "users.history":
		b = slack_history.UserHistory(*api_token, dict, flag.Args()[1])
	case "users.post":
		b = slack_post.UserPost(*api_token, dict, flag.Args()[1], string(stdin))
	case "groups.list":
		b = slack_history.PgList(dict)
	case "groups.history":
		b = slack_history.PgHistory(*api_token, dict, flag.Args()[1])
	case "groups.post":
		b = slack_post.GroupPost(*api_token, dict, flag.Args()[1], string(stdin))
	default:
		usage()
	}
	save_result(b)
}
func save_result(result string) {
	ioutil.WriteFile(*outfile, []byte(result+"\n"), os.ModePerm)
}
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s: command\n", os.Args[0])
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
`)
	os.Exit(1)
}
