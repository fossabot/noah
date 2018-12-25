package system

import (
	"io/ioutil"
	"os"
)

func CreateTempFolder() string {
	folder, err := ioutil.TempDir("", "prefix")
	if err != nil {
		panic(err)
	}
	return folder
}

func DeleteTempFolder(path string) {
	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}
}
