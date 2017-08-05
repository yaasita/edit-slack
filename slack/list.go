package slack

import (
	"fmt"
	"sort"
	"strings"
)

func (slk SlackData) ChList() string {
	var list []string
	for _, v := range slk.Channels {
		list = append(list, fmt.Sprintf("slack://ch/%s", v))
	}
	return list2str(list)
}

func (slk SlackData) PgList() string {
	var list []string
	for _, v := range slk.Groups {
		list = append(list, fmt.Sprintf("slack://pg/%s", v))
	}
	return list2str(list)
}

func (slk SlackData) UserList() string {
	var list []string
	for _, v := range slk.Users {
		list = append(list, fmt.Sprintf("slack://dm/%s", v.Name))
	}
	return list2str(list)
}

func list2str(list []string) string {
	sort.Strings(list)
	str := strings.Join(list, "\n")
	return str
}
