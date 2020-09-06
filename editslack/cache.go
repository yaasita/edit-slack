package editslack

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func (slk SlackData) SaveCache() {
	cachefile := slk.Options["cachefile"]
	slk.Options = map[string]string{}
	b, err := json.Marshal(slk)
	check(err, "json Marshal error")
	err = ioutil.WriteFile(cachefile, b, os.ModePerm)
	check(err, "write error")
}
