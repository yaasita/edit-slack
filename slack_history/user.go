package slack_history

import (
	"../slack_list"
	"fmt"
)

func UserList(ch slack_list.Top) string {
	var list []string
	for _, v := range ch.Users {
		list = append(list, fmt.Sprintf("slack://dm/%s", v.Name))
	}
	return list2str(list)
}
func UserHistory(api_token string, dict slack_list.Top, target string) string {
	if dict.User2Id(target) == "" {
		return "User Not Found"
	}
	id := dict.User2Im(api_token, target)
	raw := get_history(api_token, id, "im.history")
	messages := parse_message(raw, dict)
	return messages
}
