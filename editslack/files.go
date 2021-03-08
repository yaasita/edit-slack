package editslack

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/slack-go/slack"
)

func (slk SlackData) FileDownload(url, outfile string) {
	api := slack.New(slk.Options["token"])
	f, err1 := os.Create(outfile)
	check(err1, "os.Create error")
	err2 := api.GetFile(url, f)
	check(err2, "GetFile() error")
}
func (slk SlackData) FileUpload(file_path, chname, ts string) {
	chid := slk.name2chid(chname)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	{
		// token
		part, _ := writer.CreateFormField("token")
		part.Write([]byte(slk.Options["token"]))
	}
	{
		// file
		file, _ := os.Open(file_path)
		defer file.Close()
		part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
		io.Copy(part, file)
	}
	{
		// channels
		part, _ := writer.CreateFormField("channels")
		part.Write([]byte(chid))
	}
	if ts != "" {
		// ts
		part, _ := writer.CreateFormField("thread_ts")
		part.Write([]byte(ts))
	}
	writer.Close()
	{
		r, _ := http.NewRequest("POST", "https://slack.com/api/files.upload", body)
		r.Header.Add("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		res, err1 := client.Do(r)
		check(err1, "FileUpload error")
		if res.StatusCode != 200 {
			panic(res.Status)
		}
		type FileUploadResponse struct {
			Ok    bool   `json:"ok"`
			Error string `json:"error"`
		}
		body, err2 := ioutil.ReadAll(res.Body)
		check(err2, "ioutil.ReadAll error")
		var res_json FileUploadResponse
		err3 := json.Unmarshal(body, &res_json)
		check(err3, "Unmarshal error")
		if res_json.Ok == false {
			panic(res_json.Error)
		}
	}
}
