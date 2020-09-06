package editslack

import (
	"os"

	"github.com/slack-go/slack"
)

func (slk SlackData) FileDownload(url, outfile string) {
	api := slack.New(slk.Options["token"])
	f, err1 := os.Create(outfile)
	check(err1, "os.Create error")
	err2 := api.GetFile(url, f)
	check(err2, "GetFile() error")
}
func (slk SlackData) FileUpload(filepath, chname, ts string) {
	api := slack.New(slk.Options["token"])
	chid := slk.name2chid(chname)
	params := slack.FileUploadParameters{
		File:            filepath,
		Channels:        []string{chid},
		ThreadTimestamp: ts,
	}
	_, err := api.UploadFile(params)
	check(err, "UploadFile() error")
}
