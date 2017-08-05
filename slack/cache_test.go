package slack

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	token := os.Getenv("slack_token")
	cachefile := "/tmp/go-test"
	if stat, err := os.Stat(cachefile); err == nil && stat.Size() > 0 {
		e := os.Remove(cachefile)
		check(e, "file remove err")
		fmt.Println("FILE DELETE")
	}
	{
		slk_new := New(token, cachefile, "/tmp/go-debug")
		if fmt.Sprintf("%s", reflect.TypeOf(slk_new)) != "slack.SlackData" {
			t.Errorf("ske_new = %s", reflect.TypeOf(slk_new))
		}
		slk_new.SaveCache()
	}
	{
		slk_load := New(token, cachefile)
		if fmt.Sprintf("%s", reflect.TypeOf(slk_load)) != "slack.SlackData" {
			t.Errorf("ske_load = %s", reflect.TypeOf(slk_load))
		}
	}
}
