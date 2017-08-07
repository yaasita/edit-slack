package slack

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestHistory(t *testing.T) {
	token := os.Getenv("slack_token")
	cachefile := "/tmp/go-test"
	debugfile := "/tmp/go-debug"
	slk := New(token, cachefile, debugfile)
	var str string
	str = slk.ChHistory("server")
	ioutil.WriteFile("/tmp/go-out-ch", []byte(str), os.ModePerm)
	//str = slk.UserHistory("yamasita")
	//ioutil.WriteFile("/tmp/go-out-user", []byte(str), os.ModePerm)
	//str = slk.PgHistory("line-to-slack")
	//ioutil.WriteFile("/tmp/go-out-pg", []byte(str), os.ModePerm)
	slk.SaveCache()
}
func TestParseText(t *testing.T) {
	token := os.Getenv("slack_token")
	cachefile := "/tmp/go-test"
	debugfile := "/tmp/go-debug"
	slk := New(token, cachefile, debugfile)
	str := slk.parse_text("<@U03098XQV|hoge-huga> あ <@U174XARQ9|a.b> い <@U02V5AY5J|aba-hoge> う <@U02V5AY5J|aba_hoge>")
	fmt.Println(str)
	ioutil.WriteFile("/tmp/go-out", []byte(str), os.ModePerm)
	slk.SaveCache()
}
