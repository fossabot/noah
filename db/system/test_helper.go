package system

import (
	"io/ioutil"
	"os"
)

func CreateTempFolder() string {
	file, err := ioutil.TempFile("dir", "prefix")
	if err != nil {
		panic(err)
	}
	return file.Name()
}

func DeleteTempFolder(path string) {
	if err := os.Remove(path); err != nil {
		panic(err)
	}
}
