package editslack

import (
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
)

func (slk SlackData) RepliesHistory(chname, ts, outfile string) {
	chid := slk.name2chid(chname)
	params := slack.GetConversationRepliesParameters{
		ChannelID: chid,
		Timestamp: ts,
		Limit:     100,
	}
	api := slack.New(slk.Options["token"])
	messages, _, _, err := api.GetConversationReplies(&params)
	check(err, "GetConversationReplies error")
	var ch_messages Messages
	for _, msg := range messages {
		time_f, _ := strconv.ParseFloat(msg.Msg.Timestamp, 64)
		m := Message{
			Name:      msg.Msg.User,
			Time:      time_f,
			Text:      msg.Msg.Text,
			Files:     msg.Msg.Files,
			Reactions: msg.Msg.Reactions,
		}
		ch_messages = append(ch_messages, m)
	}
	result := slk.sprint_messages(&ch_messages)
	result += slk.add_input_line(chid)
	result += "[] Also send to channel\n"
	err2 := ioutil.WriteFile(outfile, []byte(result), os.ModePerm)
	check(err2, "write file error")
}
func (slk SlackData) RepliesPost(chname, ts, text string) {
	api := slack.New(slk.Options["token"])
	broadcast, poststr := parse_replies_text(text)
	params := slack.PostMessageParameters{
		ThreadTimestamp: ts,
		LinkNames:       1,
		ReplyBroadcast:  broadcast,
	}
	option := slack.MsgOptionPostMessageParameters(params)
	option_post := slack.MsgOptionText(poststr, true)
	chid := slk.name2chid(chname)
	_, _, err := api.PostMessage(chid, option_post, option)
	check(err, "post error")
}
func parse_replies_text(text string) (bool, string) {
	result := []string{}
	broadcast := false
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		s := scanner.Text()
		if s == "[] Also send to channel" {
			continue
		} else if s == "[x] Also send to channel" {
			broadcast = true
			continue
		}
		result = append(result, s)
	}
	err := scanner.Err()
	check(err, "scan error")
	return broadcast, strings.Join(result, "\n")
}
