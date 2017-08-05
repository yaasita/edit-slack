package slack

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestList(t *testing.T) {
	token := os.Getenv("slack_token")
	cachefile := "/tmp/go-test"
	slk := New(token, cachefile)
	str := ""
	str = slk.ChList()
	str = slk.PgList()
	str = slk.UserList()
	ioutil.WriteFile("/tmp/go-out-list", []byte(str), os.ModePerm)
	//fmt.Println(str)
}
