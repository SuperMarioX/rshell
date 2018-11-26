package update

import (
	"github.com/luckywinds/rshell/types"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var fileroot = ".rshell"
var filename = "rshell.latest"
var version types.Latestversion

func isNeedUpdate(cVersion string) int {
	return strings.Compare(cVersion, version.Version)
}

func getLatestVersion() {
	a, err := ioutil.ReadFile(fileroot + "/" + filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(a, &version)
	if err != nil {
		return
	}
}

func downloadFile(outpath string, file string, url string) error {
	out, err := os.Create(outpath + "/" + file)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url + "/" + file)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Update(c types.Cfg, cVersion string) {
	if len(c.Updateserver) == 0 {
		return
	}

	var server = ""
	for _, s := range c.Updateserver {
		downloadFile(fileroot, filename, s)
		getLatestVersion()
		if version.Version != "" && version.Release != "" {
			server = s
			break
		}
	}

	if server != "" && isNeedUpdate(cVersion) != 0 {
		downloadFile(".", version.Release, server)
	}
}