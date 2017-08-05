package slack

import (
	"fmt"
	"os"
	"testing"
)

func TestPost(t *testing.T) {
	token := os.Getenv("slack_token")
	cachefile := "/tmp/go-test"
	//slk := New(token, cachefile)
	//slk.ChannelPost("yaasita", `Debian GNU/Linux 9.1 (stretch)`)
	//slk.GroupPost("testdata", `Debian GNU/Linux 9.1 (stretch)`)
	//slk.UserPost("yamasita", `Debian GNU/Linux 9.1 (stretch)`)
	fmt.Println("ok" + token + cachefile)
}
