package slack

import (
	"fmt"
	"os"
	"testing"
)

func TestResolv(t *testing.T) {
	token := os.Getenv("slack_token")
	cachefile := "/tmp/go-test"
	slk := New(token, cachefile)
	if n, _ := slk.ID2Name("C02VD27QV"); n != "hoge" {
		t.Errorf("expect %s but got ID2Name = %s", "hoge", n)
	}
	testCases := []struct {
		testFunc string
		testStr  string
		want     string
	}{
		{"ID2Name", "C02VD27QV", "hoge"},
		{"ID2Name", "U02V5BP92", "yamasita"},
		{"ID2Name", "G61US4L5N", "line-to-slack"},
		{"ID2Name", "HOGEHOGE", ""},
		{"Channel2ID", "random", "C02V3TZ8H"},
		{"Channel2ID", "hogehogehoge", ""},
		{"Group2ID", "line-to-slack", "G61US4L5N"},
		{"Group2ID", "ABCDEFG", ""},
		{"User2ID", "yamasita", "U02V5BP92"},
		{"User2ID", "hogehogehoge", ""},
		{"User2DMID", "yamasita", "D6J7FKACA"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("function: %s input: %s want: %s", tc.testFunc, tc.testStr, tc.want), func(t *testing.T) {
			var g string
			if tc.testFunc == "ID2Name" {
				g, _ = slk.ID2Name(tc.testStr)
			} else if tc.testFunc == "Channel2ID" {
				g, _ = slk.Channel2ID(tc.testStr)
			} else if tc.testFunc == "Group2ID" {
				g, _ = slk.Group2ID(tc.testStr)
			} else if tc.testFunc == "User2ID" {
				g, _ = slk.User2ID(tc.testStr)
			} else if tc.testFunc == "User2DMID" {
				g = slk.User2DMID(tc.testStr)
			}
			if g != tc.want {
				t.Errorf("expect %s but got %s = %s", tc.testFunc, tc.want, g)
			}
		})
	}
}
