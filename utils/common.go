package utils

import (
	"github.com/golang/glog"
	"os"
)

func CheckFolder(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0777)
			if err != nil {
				glog.Fatal(err)
			}
		}
	}
}
