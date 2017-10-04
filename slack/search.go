package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

func (slk SlackData) Search(query string) string {
	var raw []byte
	data := url.Values{}
	data.Set("token", slk.Option["token"])
	data.Add("query", query)
	data.Add("count", "100")
	data.Add("sort", "timestamp")
	client := &http.Client{}
	r, err := http.NewRequest("POST", api_url+"search.messages", bytes.NewBufferString(data.Encode()))
	check(err, "request error")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(r)
	check(err, "request error client do")
	raw, err = ioutil.ReadAll(resp.Body)
	check(err, "resp body read error")
	if slk.Option["debugfile"] != "" {
		ioutil.WriteFile(slk.Option["debugfile"], raw, os.ModePerm)
	}
	return slk.parse_result(raw)
}
func (slk SlackData) parse_result(raw []byte) string {
	var messages Messages
	var json_interface interface{}
	json.Unmarshal(raw, &json_interface)
	a := json_interface.(map[string]interface{})["messages"].(map[string]interface{})["matches"].([]interface{})
	for _, v := range a {
		val := v.(map[string]interface{})
		name := ""
		time_f, _ := strconv.ParseFloat(val["ts"].(string), 64)
		text := ""
		if slk.Option["debugfile"] != "" {
			fmt.Printf("start parse type=%s,  text=%s\n", val["type"], val["text"])
		}
		if val["type"].(string) == "message" && (val["user"] != nil || val["username"] != nil) && val["text"] != nil {
			if val["user"] == nil {
				name = val["username"].(string)
			} else {
				n, _ := slk.ID2Name(val["user"].(string))
				name = n
			}
			name = "slack://ch/" + val["channel"].(map[string]interface{})["name"].(string) + " " + name
			text = slk.parse_text(val["text"].(string))
			ms := Message{
				Name:     name,
				Time:     time_f,
				Text:     text,
				File:     "",
				Reaction: "",
			}
			messages = append(messages, ms)
		}
	}
	var result string
	sort.Sort(messages)
	for _, v := range messages {
		t := time.Unix(int64(v.Time), 0)
		result += fmt.Sprintf("%s (%s): %s\n", v.Name, t.Format("2006/01/02 15:04"), v.Text)
	}
	return result
}
