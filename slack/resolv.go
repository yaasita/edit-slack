package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (slk SlackData) ID2Name(search_word string) (string, error) {
	if value, ok := slk.Channels[search_word]; ok {
		return value, nil
	}
	if value, ok := slk.Users[search_word]; ok {
		return value.Name, nil
	}
	if value, ok := slk.Groups[search_word]; ok {
		return value, nil
	}
	return "", errors.New("ID not found")
}
func (slk SlackData) Channel2ID(search_word string) (string, error) {
	for k, v := range slk.Channels {
		if v == search_word {
			return k, nil
		}
	}
	return "", errors.New("channel name not found")
}
func (slk SlackData) Group2ID(search_word string) (string, error) {
	for k, v := range slk.Groups {
		if v == search_word {
			return k, nil
		}
	}
	return "", errors.New("group name not found")
}
func (slk SlackData) User2ID(search_word string) (string, error) {
	for k, v := range slk.Users {
		if v.Name == search_word {
			return k, nil
		}
	}
	return "", errors.New("user name not found")
}

func (slk SlackData) User2DMID(user string) string {
	target_user_id, _ := slk.User2ID(user)
	dm_id := ""
	if slk.Users[target_user_id].Dm == "" {
		// api
		var raw []byte
		data := url.Values{}
		data.Set("token", slk.Option["token"])
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
		dm_id = slk.Users[target_user_id].Dm
	}
	slk.Users[target_user_id] = User{
		Name: user,
		Dm:   dm_id,
	}
	return dm_id
}
