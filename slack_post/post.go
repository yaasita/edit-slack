package slack_post

import (
	"../slack_list"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

const api_url = "https://slack.com/api/"

func ChannelPost(token string, dict slack_list.Top, target string, message string) string {
	return postch(token, dict.Channel2Id(target), message)
}
func UserPost(token string, dict slack_list.Top, target string, message string) string {
	return postch(token, dict.User2Im(token, target), message)
}
func GroupPost(token string, dict slack_list.Top, target string, message string) string {
	return postch(token, dict.Group2Id(target), message)
}
func postch(token, id, message string) string {
	data := url.Values{}
	data.Set("token", token)
	data.Add("channel", id)
	data.Add("text", message)
	data.Add("as_user", "true")
	data.Add("link_names", "1")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url+"chat.postMessage", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(r)
	raw, _ := ioutil.ReadAll(resp.Body)
	return string(raw)
}
