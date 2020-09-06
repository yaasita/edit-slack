package editslack

import (
	"fmt"
	"testing"
)

func TestConversationJoin(t *testing.T) {
	slk.ConversationJoin("ch-example")
	slk.ConversationJoin("ch-example2")
	slk.ConversationJoin("ch-example3")
	slk.SaveCache()
}

func TestConversationLeave(t *testing.T) {
	c := slk.name2chid("mpdm-spam4--spam2--yaasita-1")
	fmt.Println(c)
	slk.ConversationLeave("mpdm-spam4--spam2--yaasita-1")
	slk.SaveCache()
}

func TestConversationPost(t *testing.T) {
	slk.ConversationPost("general", "hogehoge")
}
