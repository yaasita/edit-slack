package slack

import (
	"bytes"
	"net/http"
	"net/url"
)

func (slk SlackData) ChannelPost(name string, message string) {
	id, err := slk.Channel2ID(name)
	check(err, "channel post: ch not found")
	slk.postch(id, message)
}
func (slk SlackData) GroupPost(name string, message string) {
	id, err := slk.Group2ID(name)
	check(err, "group post: ch not found")
	slk.postch(id, message)
}
func (slk SlackData) UserPost(name string, message string) {
	id := slk.User2DMID(name)
	slk.postch(id, message)
}
func (slk SlackData) postch(id, message string) {
	data := url.Values{}
	data.Set("token", slk.Option["token"])
	data.Add("channel", id)
	data.Add("text", message)
	data.Add("as_user", "true")
	data.Add("link_names", "1")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url+"chat.postMessage", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, _ = client.Do(r)
}
