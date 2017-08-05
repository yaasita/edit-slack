package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const api_url = "https://slack.com/api/"

type User struct {
	Name string `json:"name"`
	Dm   string `json:dm`
}
type SlackData struct {
	Channels map[string]string `json:"channels"`
	Users    map[string]User   `json:"users"`
	Groups   map[string]string `json:"groups"`
	Option   map[string]string `json:"option"`
}

func (slk SlackData) SaveCache() {
	b, err := json.Marshal(slk)
	check(err, "json Marshal error")
	ioutil.WriteFile(slk.Option["cachefile"], b, os.ModePerm)
}

func New(token, cachefile string, debugfile ...string) SlackData {
	slack_data := SlackData{}
	if stat, err := os.Stat(cachefile); err == nil && stat.Size() > 0 {
		slack_data = load_cache(cachefile)
	} else {
		slack_data = create_cache(cachefile, token)
	}
	// option
	slack_data.Option = map[string]string{}
	slack_data.Option["cachefile"] = cachefile
	slack_data.Option["token"] = token
	if len(debugfile) == 1 {
		slack_data.Option["debugfile"] = debugfile[0]
	} else {
		slack_data.Option["debugfile"] = ""
	}
	return slack_data
}

func load_cache(cachefile string) SlackData {
	slack_data := SlackData{}
	data, err := ioutil.ReadFile(cachefile)
	check(err, "load cachefile error")
	json.Unmarshal(data, &slack_data)
	return slack_data
}
func create_cache(cachefile, token string) SlackData {
	slack_data := SlackData{}
	{
		// channel
		slack_data.Channels = map[string]string{}
		var raw []byte
		data := url.Values{}
		data.Set("token", token)
		data.Add("exclude_archive", "1")
		client := &http.Client{}
		r, _ := http.NewRequest("POST", api_url+"channels.list", bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, _ := client.Do(r)
		raw, _ = ioutil.ReadAll(resp.Body)
		var json_interface interface{}
		json.Unmarshal(raw, &json_interface)
		a := json_interface.(map[string]interface{})
		b := a["channels"]
		c := b.([]interface{})
		for _, v := range c {
			val := v.(map[string]interface{})
			slack_data.Channels[val["id"].(string)] = val["name"].(string)
		}
	}
	{
		// user
		slack_data.Users = map[string]User{}
		var raw []byte
		data := url.Values{}
		data.Set("token", token)
		client := &http.Client{}
		r, _ := http.NewRequest("POST", api_url+"users.list", bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, _ := client.Do(r)
		raw, _ = ioutil.ReadAll(resp.Body)
		var json_interface interface{}
		json.Unmarshal(raw, &json_interface)
		a := json_interface.(map[string]interface{})
		b := a["members"]
		c := b.([]interface{})
		for _, v := range c {
			val := v.(map[string]interface{})
			if val["deleted"].(bool) == true {
				continue
			}
			user := User{
				Name: val["name"].(string),
				Dm:   "",
			}
			slack_data.Users[val["id"].(string)] = user
		}
	}
	{
		// private group
		slack_data.Groups = map[string]string{}
		var raw []byte
		data := url.Values{}
		data.Set("token", token)
		client := &http.Client{}
		r, _ := http.NewRequest("POST", api_url+"groups.list", bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, _ := client.Do(r)
		raw, _ = ioutil.ReadAll(resp.Body)
		var json_interface interface{}
		json.Unmarshal(raw, &json_interface)
		a := json_interface.(map[string]interface{})
		b := a["groups"]
		c := b.([]interface{})
		for _, v := range c {
			val := v.(map[string]interface{})
			slack_data.Groups[val["id"].(string)] = val["name"].(string)
		}
	}
	return slack_data
}
