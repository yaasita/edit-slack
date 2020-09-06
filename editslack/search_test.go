package editslack

import (
	"testing"
)

func TestSearchAll(t *testing.T) {
	slk.SearchAll("a b", "/tmp/o.txt")
}
