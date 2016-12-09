package slack_list

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const api_url = "https://slack.com/api/"

type Channel struct {
	Name string `json:"name"`
}
type Group Channel
type User struct {
	Name string `json:"name"`
	Dm   string `json:dm`
}
type Option struct {
	Cachefile string `json:"cachefile"`
}
type Top struct {
	Channels map[string]Channel `json:"channels"`
	Users    map[string]User    `json:"users"`
	Groups   map[string]Group   `json:"groups"`
	Option   Option             `json:"option"`
}

func (t Top) Id2Name(search_word string) string {
	if value, ok := t.Channels[search_word]; ok {
		return value.Name
	}
	if value, ok := t.Users[search_word]; ok {
		return value.Name
	}
	if value, ok := t.Groups[search_word]; ok {
		return value.Name
	}
	return ""
}
func (t Top) Channel2Id(search_word string) string {
	for k, v := range t.Channels {
		if v.Name == search_word {
			return k
		}
	}
	return ""
}
func (t Top) Group2Id(search_word string) string {
	for k, v := range t.Groups {
		if v.Name == search_word {
			return k
		}
	}
	return ""
}
func (t Top) User2Id(search_word string) string {
	for k, v := range t.Users {
		if v.Name == search_word {
			return k
		}
	}
	return ""
}
func (t Top) User2Im(api_token, user string) string {
	var target_user_id = t.User2Id(user)
	var dm_id string
	if t.Users[target_user_id].Dm == "" {
		// api
		var raw []byte
		data := url.Values{}
		data.Set("token", api_token)
		data.Add("user", target_user_id)
		client := &http.Client{}
		r, _ := http.NewRequest("POST", api_url+"im.open", bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, _ := client.Do(r)
		raw, _ = ioutil.ReadAll(resp.Body)
		var json_interface interface{}
		json.Unmarshal(raw, &json_interface)
		a := json_interface.(map[string]interface{})
		b := a["channel"]
		c := b.(map[string]interface{})
		d := c["id"]
		dm_id = d.(string)
	} else {
		dm_id = t.Users[target_user_id].Dm
	}
	t.Users[target_user_id] = User{
		Name: user,
		Dm:   dm_id,
	}
	save_hash(t, t.Option.Cachefile)
	return dm_id
}

func List(token string, cachefile string) Top {
	if stat, err := os.Stat(cachefile); err != nil || stat.Size() == 0 {
		create_hash(token, cachefile)
	}
	t := load_hash(cachefile)
	return t
}
func create_hash(token, cachefile string) {
	top := Top{}
	{
		// channel
		top.Channels = map[string]Channel{}
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
			top.Channels[val["id"].(string)] = Channel{
				Name: val["name"].(string),
			}
		}
	}
	{
		// user
		top.Users = map[string]User{}
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
			top.Users[val["id"].(string)] = user
		}
	}
	{
		// private group
		top.Groups = map[string]Group{}
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
			pg := Group{
				Name: val["name"].(string),
			}
			top.Groups[val["id"].(string)] = pg
		}
	}
	top.Option.Cachefile = cachefile
	save_hash(top, cachefile)
}
func save_hash(hashdata Top, cachefile string) {
	b, _ := json.Marshal(hashdata)
	ioutil.WriteFile(cachefile, b, os.ModePerm)
}
func load_hash(load_file string) Top {
	data, err := ioutil.ReadFile(load_file)
	if err != nil {
		panic(err)
	}
	top := Top{}
	json.Unmarshal(data, &top)
	return top
}
