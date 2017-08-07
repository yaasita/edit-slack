package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Name     string
	Time     float64
	Text     string
	File     string
	Reaction string
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

func (slk SlackData) ChHistory(chname string) string {
	return slk.get_history(chname, "channels.history")
}
func (slk SlackData) UserHistory(uname string) string {
	return slk.get_history(uname, "im.history")
}
func (slk SlackData) PgHistory(pgname string) string {
	return slk.get_history(pgname, "groups.history")
}
func (slk SlackData) get_history(name, api_kind string) string {
	var id string
	{
		var err error
		if api_kind == "channels.history" {
			id, err = slk.Channel2ID(name)
			check(err, "ChHistory channel not found")
		} else if api_kind == "im.history" {
			id = slk.User2DMID(name)
			check(err, "UserHistory channel not found")
		} else if api_kind == "groups.history" {
			id, err = slk.Group2ID(name)
			check(err, "Group name not found")
		}
	}
	// get api
	var raw []byte
	data := url.Values{}
	data.Set("token", slk.Option["token"])
	data.Add("count", "300")
	data.Add("channel", id)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url+api_kind, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(r)
	raw, _ = ioutil.ReadAll(resp.Body)
	if slk.Option["debugfile"] != "" {
		ioutil.WriteFile(slk.Option["debugfile"], raw, os.ModePerm)
	}
	return slk.parse_message(raw)
}
func (slk SlackData) parse_message(raw []byte) string {
	var ch_messages Messages
	var json_interface interface{}
	json.Unmarshal(raw, &json_interface)
	a := json_interface.(map[string]interface{})
	b := a["messages"]
	c := b.([]interface{})
	for _, v := range c {
		val := v.(map[string]interface{})
		name := ""
		time_f, _ := strconv.ParseFloat(val["ts"].(string), 64)
		file := ""
		text := ""
		reaction := ""
		if val["type"].(string) == "message" && (val["user"] != nil || val["username"] != nil) && val["text"] != nil {
			//fmt.Println(val["ts"].(string))
			if val["user"] == nil {
				name = val["username"].(string)
			} else {
				n, _ := slk.ID2Name(val["user"].(string))
				name = n
			}
			if f := val["file"]; f != nil {
				fn := f.(map[string]interface{})["url_private"].(string)
				file = fmt.Sprintf("=>download link: %s\n", fn)
			}
			text = slk.parse_text(val["text"].(string))
			if r := val["reactions"]; r != nil {
				reaction = slk.parse_reaction(r)
			}
			ms := Message{
				Name:     name,
				Time:     time_f,
				Text:     text,
				File:     file,
				Reaction: reaction,
			}
			ch_messages = append(ch_messages, ms)
		}
	}
	var result string
	sort.Sort(ch_messages)
	for _, v := range ch_messages {
		t := time.Unix(int64(v.Time), 0)
		result += fmt.Sprintf("%s (%s): %s\n", v.Name, t.Format("2006/01/02 15:04"), v.Text)
		result += v.File
		result += v.Reaction
	}
	return result
}
func (slk SlackData) parse_reaction(reaction interface{}) string {
	var result string
	a := reaction.([]interface{})
	for _, v := range a {
		result += "=>reaction: "
		b := v.(map[string]interface{})
		result += ":" + b["name"].(string) + ": "
		c := b["users"]
		d := c.([]interface{})
		var users []string
		for _, v2 := range d {
			name, _ := slk.ID2Name(v2.(string))
			users = append(users, "@"+name)
		}
		result += strings.Join(users, ", ")
		result += "\n"
	}
	return result
}
func (slk SlackData) parse_text(intext string) string {
	var result string
	{
		// user/ch name replace
		r := regexp.MustCompile(`<[@#][\w\|\-\.]+>`)
		result = r.ReplaceAllStringFunc(intext, func(match string) string {
			re := regexp.MustCompile(`<([@#])([\w\-\.]+)`)
			m := re.FindStringSubmatch(match)
			id, _ := slk.ID2Name(m[2])
			return m[1] + id
		})
	}
	{
		// remove < >
		r := regexp.MustCompile(`[<>]`)
		result = r.ReplaceAllString(result, " ")
	}
	{
		// replace &gt; &lt;
		r1 := regexp.MustCompile(`&gt;`)
		result = r1.ReplaceAllString(result, ">")
		r2 := regexp.MustCompile(`&lt;`)
		result = r2.ReplaceAllString(result, "<")
	}
	return result
}
