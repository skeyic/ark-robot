package utils

import (
	"flag"
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

func EnableGlogForTesting() {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()
}
