package slack

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	token := os.Getenv("slack_token")
	cachefile := "/tmp/go-test"
	debugfile := "/tmp/go-debug"
	slk := New(token, cachefile, debugfile)
	str := slk.Search("from:@yamasita")
	ioutil.WriteFile("/tmp/go-out-search", []byte(str+"\n"), os.ModePerm)
	slk.SaveCache()
}

//func TestParse(t *testing.T) {
//	token := os.Getenv("slack_token")
//	cachefile := "/tmp/go-test"
//	debugfile := "/tmp/go-debug"
//	slk := New(token, cachefile, debugfile)
//	in, _ := ioutil.ReadFile("/tmp/go-debug.json")
//	str := slk.parse_result(in)
//	ioutil.WriteFile("/tmp/go-out-search", []byte(str), os.ModePerm)
//	slk.SaveCache()
//}
