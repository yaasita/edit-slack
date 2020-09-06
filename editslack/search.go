package editslack

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/slack-go/slack"
)

func (slk SlackData) SearchAll(search_word, outfile string) {
	api := slack.New(slk.Options["token"])
	params := slack.SearchParameters{
		Sort:          "score",
		SortDirection: "desc",
		Highlight:     false,
		Count:         100,
		Page:          1,
	}
	messages, files, err1 := api.Search(search_word, params)
	check(err1, "search error")
	result := ""
	{
		result += "# messages\n"
		var ch_messages Messages
		for _, m := range messages.Matches {
			time_f, _ := strconv.ParseFloat(m.Timestamp, 64)
			chid := m.Channel.ID
			m := Message{
				Uri:  fmt.Sprintf("slack://ch/%s", slk.Channels[chid].Name),
				Name: m.User,
				Time: time_f,
				Text: m.Text,
			}
			ch_messages = append(ch_messages, m)
		}
		result += slk.sprint_messages(&ch_messages)
		result += "\n"
	}
	{
		result += "# files\n"
		var ch_messages Messages
		for _, f := range files.Matches {
			time_f := float64(f.Timestamp)
			if len(f.Channels) == 0 {
				continue
			}
			chid := f.Channels[0]
			m := Message{
				Uri:   fmt.Sprintf("slack://ch/%s", slk.Channels[chid].Name),
				Name:  f.User,
				Time:  time_f,
				Files: []slack.File{f},
			}
			ch_messages = append(ch_messages, m)
		}
		result += slk.sprint_messages(&ch_messages)
	}
	err2 := ioutil.WriteFile(outfile, []byte(result), os.ModePerm)
	check(err2, "write file error")
}
