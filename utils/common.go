package utils

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
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

func CheckFile(path string) {
	_, err := os.Stat(path)
	if err != nil {
		glog.Fatal(err)
	}
}

func CheckFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		glog.Fatal(err)
	}
}

func CopyFile(sourceFile, destinationFile string) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		glog.Errorf("Error reading %s, err: %v", sourceFile, err)
		glog.Fatal(err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		glog.Errorf("Error creating %s, err: %v", destinationFile, err)
		glog.Fatal(err)
		return
	}
}

func EnableGlogForTesting() {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()
}

func ThousandFormatFloat64(num float64) string {
	numStr := fmt.Sprintf("%.0f", num)
	length := len(numStr)
	if length < 4 {
		return numStr
	}
	count := (length - 1) / 3
	for i := 0; i < count; i++ {
		numStr = numStr[:length-(i+1)*3] + "," + numStr[length-(i+1)*3:]
	}
	return numStr
}
