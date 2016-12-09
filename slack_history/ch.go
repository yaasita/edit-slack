package slack_history

import (
	"../slack_list"
	"fmt"
)

func ChList(ch slack_list.Top) string {
	var list []string
	for _, v := range ch.Channels {
		list = append(list, fmt.Sprintf("slack://ch/%s", v.Name))
	}
	return list2str(list)
}
func ChHistory(api_token string, dict slack_list.Top, target string) string {
	if dict.Channel2Id(target) == "" {
		return "Channel Not Found"
	}
	raw := get_history(api_token, dict.Channel2Id(target), "channels.history")
	messages := parse_message(raw, dict)
	return messages
}
