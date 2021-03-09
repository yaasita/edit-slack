package editslack

import (
	"os"
	"testing"
)

var token = os.Getenv("MY_SLACK_TOKEN")
var cachefile = "/tmp/c.json"
var slk = New(token, cachefile)

func TestName2chid(t *testing.T) {
	name2id := func(name string) {
		id := slk.name2chid(name)
		t.Logf("get id = %s", id)
		if id == "" {
			t.Error("id is blank")
		}
	}
	name2id("ch-example")
	//name2id("general")
	//name2id("ch-example2")
	//name2id("ch-example3")
}
func TestFormatTextCh(t *testing.T) {
	ch := "<#C3HJSM2CX|ch-example> <#C3HJSPGK1|ch-example2> <#C01N4D27W79|under_bar>"
	after := slk.format_text(ch)
	//t.Logf("%s\n", after)
	if after != "#ch-example #ch-example2 #under_bar\n" {
		t.Errorf("return = %s\n", after)
	}
}
func TestFormatTextUser(t *testing.T) {
	user := "<@U038TDF1P> <@U038TFMRB> <@U038UGF8N>"
	after := slk.format_text(user)
	t.Logf("%s\n", after)
	//if after != "#ch-example #ch-example2 #under_bar\n" {
	//	t.Errorf("return = %s\n", after)
	//}
}
func TestFormatTextEscape(t *testing.T) {
	before := "&amp;abc"
	after := slk.format_text(before)
	want := "&abc\n"
	if after != want {
		t.Errorf("want = %#v", want)
		t.Errorf("after = %#v", after)
		t.Errorf("before = %#v", before)
	}
}
