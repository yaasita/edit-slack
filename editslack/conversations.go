package editslack

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
)

func (slk SlackData) ConversationList(outfile string) {
	var list string
	{
		// channels
		channels := []Channel{}
		for _, v := range slk.Channels {
			channels = append(channels, v)
		}
		sort.Slice(channels, func(i, j int) bool {
			return channels[i].Name < channels[j].Name
		})
		list = "# public\n"
		for _, v := range channels {
			if !v.IsIM && !v.IsPrivate && v.IsMember {
				list += fmt.Sprintf("slack://ch/%s\n", v.Name)
			}
		}

		list += "\n# public (not in member)\n"
		for _, v := range channels {
			if !v.IsIM && !v.IsPrivate && !v.IsMember {
				list += fmt.Sprintf("slack://ch/%s\n", v.Name)
			}
		}

		list += "\n# private\n"
		for _, v := range channels {
			if v.IsPrivate {
				list += fmt.Sprintf("slack://ch/%s\n", v.Name)
			}
		}
	}
	{
		// users
		list += "\n# dm\n"
		users := []User{}
		for _, v := range slk.Users {
			users = append(users, v)
		}
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name < users[j].Name
		})
		for _, v := range users {
			list += fmt.Sprintf("slack://ch/%s (%s)\n", v.Name, v.DisplayName)
		}
	}
	err := ioutil.WriteFile(outfile, []byte(list), os.ModePerm)
	check(err, "write file error")
}
func (slk SlackData) ConversationHistory(chname, outfile string) {
	chid := slk.name2chid(chname)
	api := slack.New(slk.Options["token"])
	params := slack.GetConversationHistoryParameters{
		ChannelID: chid,
		Inclusive: false,
		Limit:     100,
	}
	var ch_messages Messages
	response, err1 := api.GetConversationHistory(&params)
	check(err1, "GetConversationHistory error")
	for _, msg := range response.Messages {
		time_f, _ := strconv.ParseFloat(msg.Timestamp, 64)
		m := Message{
			Uri:        fmt.Sprintf("slack://ch/%s/%s", chname, msg.Timestamp),
			Name:       msg.User,
			Time:       time_f,
			Text:       msg.Text,
			Files:      msg.Files,
			Reactions:  msg.Reactions,
			RepliesURI: "",
		}
		if msg.ReplyCount > 0 {
			m.RepliesURI = fmt.Sprintf("=> replies (%d): slack://ch/%s/%s\n",
				msg.ReplyCount,
				chname,
				msg.Timestamp)
		}
		ch_messages = append(ch_messages, m)
	}
	result := slk.sprint_messages(&ch_messages)
	result += slk.add_input_line(chid)
	err2 := ioutil.WriteFile(outfile, []byte(result), os.ModePerm)
	check(err2, "write file error")
}
func (slk SlackData) ConversationPost(chname, text string) {
	api := slack.New(slk.Options["token"])
	params := slack.PostMessageParameters{
		LinkNames: 1,
	}
	option := slack.MsgOptionPostMessageParameters(params)
	option_post := slack.MsgOptionText(strings.TrimRight(text, "\n"), true)
	chid := slk.name2chid(chname)
	_, _, err := api.PostMessage(chid, option_post, option)
	check(err, "post error")
}
func (slk *SlackData) ConversationJoin(chname string) {
	api := slack.New(slk.Options["token"])
	chid := slk.name2chid(chname)
	_, _, _, err := api.JoinConversation(chid)
	check(err, "join error")
	ch_conf := slk.Channels[chid]
	ch_conf.IsMember = true
	slk.Channels[chid] = ch_conf
}
func (slk *SlackData) ConversationLeave(chname string) {
	api := slack.New(slk.Options["token"])
	chid := slk.name2chid(chname)
	_, err := api.LeaveConversation(chid)
	check(err, "leave error")
	if slk.Channels[chid].IsPrivate {
		delete(slk.Channels, chid)
	} else {
		ch_conf := slk.Channels[chid]
		ch_conf.IsMember = false
		slk.Channels[chid] = ch_conf
	}
}
