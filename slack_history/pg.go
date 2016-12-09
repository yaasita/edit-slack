package slack_history

import (
	"../slack_list"
	"fmt"
)

func PgList(ch slack_list.Top) string {
	var list []string
	for _, v := range ch.Groups {
		list = append(list, fmt.Sprintf("slack://pg/%s", v.Name))
	}
	return list2str(list)
}
func PgHistory(api_token string, dict slack_list.Top, target string) string {
	if dict.Group2Id(target) == "" {
		return "Channel Not Found"
	}
	raw := get_history(api_token, dict.Group2Id(target), "groups.history")
	messages := parse_message(raw, dict)
	return messages
}
