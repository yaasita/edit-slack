package editslack

import (
	"strings"
	"testing"
)

func TestRepliesPost(t *testing.T) {
	slk.RepliesPost("ch-example2", "1612807446.000300",
		`[] Also send to channel
hello
foo
bar`)
	slk.RepliesPost("ch-example2", "1612807446.000300",
		`[x] Also send to channel
broad
cast
message!
`)
}
func TestParseRepliesText(t *testing.T) {
	b1, s1 := parse_replies_text(
		`aaaaaaaaaaaaa
bbbbbbbbb
[] Also send to channel
abc
def
`)
	if b1 != false {
		t.Error("b1 is not broadcast message")
	}
	if strings.Index(s1, "Also send to channel") != -1 {
		t.Errorf("%s", s1)
	}
	b2, s2 := parse_replies_text(
		`eeeeeeeeeeeeeeeee
fffffffffffff
[x] Also send to channel

abc
def
`)
	if b2 != true {
		t.Error("b2 is broadcast message")
	}
	if strings.Index(s2, "Also send to channel") != -1 {
		t.Errorf("%s", s2)
	}
}
