package editslack

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/slack-go/slack"
)

type Channel struct {
	IsMember  bool   `json:"member"`
	IsPrivate bool   `json:"private"`
	IsIM      bool   `json:"im"`
	IsMpIM    bool   `json:"mpim"`
	Name      string `json:"name"`
}
type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
type SlackData struct {
	Channels map[string]Channel `json:"channels"`
	Users    map[string]User    `json:"users"`
	Options  map[string]string
}

func load_cache(cachefile string) SlackData {
	var slackdata SlackData
	data, err := ioutil.ReadFile(cachefile)
	check(err, "load cachefile error")
	json.Unmarshal(data, &slackdata)
	return slackdata
}
func create_cache_data(token string) SlackData {
	var slackdata SlackData
	slackdata.Channels = map[string]Channel{}
	slackdata.Users = map[string]User{}
	api := slack.New(token)
	{
		// users list
		users, err := api.GetUsers()
		check(err, "error GetUsers()")
		for _, u := range users {
			if !u.Deleted {
				var display_name string
				if u.RealName != "" {
					display_name = u.RealName
				}
				if u.Profile.DisplayNameNormalized != "" {
					display_name = u.Profile.DisplayNameNormalized
				}
				slackdata.Users[u.ID] = User{Name: u.Name, DisplayName: display_name}
			}
		}
	}
	{
		// channels list
		params := slack.GetConversationsParameters{
			Cursor:          "",
			ExcludeArchived: "true",
			Limit:           100,
			Types:           []string{"public_channel", "private_channel", "mpim", "im"},
		}
		for {
			ch, nextCursor, err := api.GetConversations(&params)
			check(err, "error GetConversations")
			for _, c := range ch {
				id := c.GroupConversation.Conversation.ID
				name := c.GroupConversation.Name
				is_member := c.IsMember
				is_private := c.GroupConversation.Conversation.IsPrivate
				is_mpim := c.GroupConversation.Conversation.IsMpIM
				if c.GroupConversation.Conversation.IsIM {
					user_id := c.GroupConversation.Conversation.User
					user_name := slackdata.Users[user_id].Name
					if user_name != "" {
						slackdata.Channels[id] = Channel{Name: user_name, IsIM: true}
					}
				} else {
					slackdata.Channels[id] = Channel{Name: name, IsMember: is_member, IsPrivate: is_private, IsMpIM: is_mpim}
				}
			}
			if nextCursor == "" {
				break
			}
			params.Cursor = nextCursor
		}
	}
	return slackdata
}
func New(token, cachefile string) SlackData {
	var slackdata SlackData
	if stat, err := os.Stat(cachefile); err == nil && stat.Size() > 0 {
		slackdata = load_cache(cachefile)
	} else {
		slackdata = create_cache_data(token)
	}
	slackdata.Options = map[string]string{"token": token, "cachefile": cachefile}
	return slackdata
}
