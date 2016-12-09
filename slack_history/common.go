package slack_history

import (
	"../slack_list"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const api_url = "https://slack.com/api/"

type Message struct {
	Name string
	Time float64
	Text string
}
type Messages []Message

func (m Messages) Len() int {
	return len(m)
}
func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m Messages) Less(i, j int) bool {
	return m[i].Time < m[j].Time
}

func list2str(list []string) string {
	sort.Strings(list)
	str := strings.Join(list, "\n")
	return str
}
func get_history(api_token string, id string, api_kind string) []byte {
	// get api
	var raw []byte
	data := url.Values{}
	data.Set("token", api_token)
	data.Add("count", "300")
	data.Add("channel", id)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url+api_kind, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(r)
	raw, _ = ioutil.ReadAll(resp.Body)
	return raw
}
func parse_message(raw []byte, dict slack_list.Top) string {
	var ch_messages Messages
	var json_interface interface{}
	json.Unmarshal(raw, &json_interface)
	a := json_interface.(map[string]interface{})
	b := a["messages"]
	c := b.([]interface{})
	for _, v := range c {
		val := v.(map[string]interface{})
		time_f, _ := strconv.ParseFloat(val["ts"].(string), 64)
		if val["type"].(string) == "message" && val["user"] != nil {
			ms := Message{
				Name: val["user"].(string),
				Time: time_f,
				Text: val["text"].(string),
			}
			ch_messages = append(ch_messages, ms)
		}
	}
	var messages string
	sort.Sort(ch_messages)
	for _, v := range ch_messages {
		t := time.Unix(int64(v.Time), 0)
		messages = messages + fmt.Sprintf("%s (%s): %s\n", dict.Id2Name(v.Name), t.Format("2006/01/02 15:04"), v.Text)
	}
	return messages
}
