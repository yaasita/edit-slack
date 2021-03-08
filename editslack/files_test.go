package editslack

import (
	"testing"
)

func TestChFileUpload(t *testing.T) {
	slk.FileUpload("/tmp/a.txt", "general", "")
}
func TestRpFileUpload(t *testing.T) {
	slk.FileUpload("/tmp/a.txt", "general", "1615217889.001100")
}
