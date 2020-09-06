package editslack

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

type Message struct {
	Uri        string
	Name       string
	Time       float64
	Text       string
	Files      []slack.File
	Reactions  []slack.ItemReaction
	RepliesURI string
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

func check(e error, s string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, s)
		panic(e)
	}
}
func (slk SlackData) name2chid(name string) string {
	for k, v := range slk.Channels {
		if v.Name == name {
			return k
		}
	}
	return ""
}
func (slk SlackData) format_reactions(reactions []slack.ItemReaction) string {
	var result string
	for _, val := range reactions {
		var users []string
		for _, u := range val.Users {
			users = append(users, slk.Users[u].DisplayName)
		}
		result += fmt.Sprintf("=> reaction %s (%d): %s\n", val.Name, len(users), strings.Join(users, ", "))
	}
	return result
}
func (slk SlackData) format_text(msg string) string {
	var result string
	{
		msg_a := strings.Split(msg, "\n")
		for index, val := range msg_a {
			if index == 0 {
				result += fmt.Sprintf("%s\n", val)
			} else {
				result += fmt.Sprintf("\t%s\n", val)
			}
		}
	}
	for {
		// ch name replace
		r := regexp.MustCompile(`<#(\w+)\|([\w\-\_]+)>`)
		m := r.FindStringSubmatch(result)
		if len(m) > 0 {
			result = strings.Replace(result, m[0], "#"+m[2], -1)
		} else {
			break
		}
	}
	for {
		// user name replace
		r := regexp.MustCompile(`<@(\w+)>`)
		m := r.FindStringSubmatch(result)
		if len(m) > 0 {
			result = strings.Replace(result, m[0], "@"+slk.Users[m[1]].DisplayName, -1)
		} else {
			break
		}
	}
	// replace &gt; &lt;
	result = strings.Replace(result, `&gt;`, ">", -1)
	result = strings.Replace(result, `&lt;`, "<", -1)
	return result
}
func format_files(files []slack.File) string {
	var result string
	for _, f := range files {
		result += fmt.Sprintf("=> file: %s\n", f.URLPrivateDownload)
	}
	return result
}
func (slk SlackData) sprint_messages(ch_messages *Messages) string {
	result := ""
	sort.Sort(*ch_messages)
	for _, m := range *ch_messages {
		result += slk.format_message(m.Uri,
			m.Name,
			strconv.FormatFloat(m.Time, 'f', -1, 64),
			m.Text)
		result += slk.format_reactions(m.Reactions)
		result += format_files(m.Files)
		result += fmt.Sprintf("%s", m.RepliesURI)
	}
	return result
}
func (slk SlackData) add_input_line(chid string) string {
	if slk.Channels[chid].IsMember || slk.Channels[chid].IsIM {
		return "=== Message ===\n"
	} else {
		return "=== Not in member ===\n"
	}
}
func (slk SlackData) format_message(uri, name, timestamp, text string) string {
	s1 := ""
	if uri != "" {
		s1 = uri + " "
	}

	s2 := slk.Users[name].DisplayName

	s3 := ""
	{
		r1 := regexp.MustCompile(`\.\d+$`)
		num := r1.ReplaceAllString(timestamp, "")
		unixtime, err := strconv.ParseInt(num, 10, 64)
		check(err, "ParseInt error")
		t := time.Unix(unixtime, 0)
		s3 = t.Format("2006/01/02 15:04")
	}

	s4 := slk.format_text(text)

	return fmt.Sprintf("%s%s (%s): %s", s1, s2, s3, s4)
}
