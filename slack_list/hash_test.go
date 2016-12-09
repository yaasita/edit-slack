package slack_list

import (
	"os"
	"testing"
)

var Config map[string]string = map[string]string{}

func TestMain(m *testing.M) {
	Config["token"] = os.Getenv("slack_token")
	Config["cachefile"] = "/tmp/go-test"
	code := m.Run()
	//os.Remove(Config["cachefile"])
	os.Exit(code)
}

func TestId2Name(t *testing.T) {
	dict := List(Config["token"], Config["cachefile"])
	id := dict.Id2Name("C038TDF2B")
	if id != "general" {
		t.Errorf("Id2Name = %s", id)
	}
}
func TestChannel2Id(t *testing.T) {
	dict := List(Config["token"], Config["cachefile"])
	id := dict.Channel2Id("general")
	if id != "C038TDF2B" {
		t.Errorf("Channel2Id = %s", id)
	}
}
func TestGroup2Id(t *testing.T) {
	dict := List(Config["token"], Config["cachefile"])
	id := dict.Group2Id("aab")
	if id != "G3C9B6P7U" {
		t.Errorf("Group2Id = %s", id)
	}
}
func TestUser2Id(t *testing.T) {
	dict := List(Config["token"], Config["cachefile"])
	id := dict.User2Id("yaasita")
	if id != "U038TDF1P" {
		t.Errorf("User2Id = %s", id)
	}
}
func TestUser2Im(t *testing.T) {
	dict := List(Config["token"], Config["cachefile"])
	id := dict.User2Im(Config["token"], "yaasita")
	if id != "D3CBHQ3NY" {
		t.Errorf("DM ID = %s", id)
	}
}
