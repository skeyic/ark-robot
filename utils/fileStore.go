package utils

import (
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func WriteAndSyncFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err == nil {
		err = f.Sync()
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

type FileStoreSvc struct {
	fileName string
}

func NewFileStoreSvc(fileName string) *FileStoreSvc {
	return &FileStoreSvc{
		fileName: fileName,
	}
}

func (f *FileStoreSvc) Save(content []byte) error {
	return SaveToFile(f.fileName, content)
}

func (f *FileStoreSvc) Read() (content []byte, err error) {
	return ReadFromFile(f.fileName)
}

func SaveToFile(fileName string, content []byte) (err error) {
	//err = WriteAndSyncFile(fileName, content, 0666)
	return ioutil.WriteFile(fileName, content, 0666)
}

func ReadFromFile(fileName string) (content []byte, err error) {
	return ioutil.ReadFile(fileName)
}

type MultiFileStoreSvc struct {
	folderPath string
	prefix     string
}

func NewMultiFileStoreSvc(folderPath, prefix string) *MultiFileStoreSvc {
	_, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(folderPath, 0777)
			if err != nil {
				panic(err)
			}
		}
	}
	_, err = os.Stat(folderPath + "/" + prefix)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(folderPath+"/"+prefix, 0777)
			if err != nil {
				panic(err)
			}
		}
	}
	return &MultiFileStoreSvc{
		folderPath: folderPath,
		prefix:     prefix,
	}
}

func (f *MultiFileStoreSvc) ToRealFileName(fileName string) string {
	return f.folderPath + "/" + f.prefix + fileName
}

func (f *MultiFileStoreSvc) ToFileName(fileName string) string {
	return strings.TrimPrefix(fileName, f.folderPath+"/"+f.prefix)
}

func (f *MultiFileStoreSvc) Save(fileName string, content []byte) error {
	return SaveToFile(f.ToRealFileName(fileName), content)
}

func (f *MultiFileStoreSvc) Read(fileName string) (content []byte, err error) {
	return ReadFromFile(f.ToRealFileName(fileName))
}

type FileContent struct {
	Name    string
	Content []byte
}

func (f *MultiFileStoreSvc) ReadAll() (fileContents []*FileContent, err error) {
	var (
		files []string
	)

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), f.prefix) {
			files = append(files, strings.TrimPrefix(info.Name(), f.prefix))
		}
		return nil
	}
	err = filepath.Walk(f.folderPath, walkFunc)
	if err != nil {
		return
	}

	for _, theFile := range files {
		content, err := f.Read(theFile)
		if err != nil {
			glog.Error(err)
			panic("failed to read " + theFile)
		}
		fileContents = append(fileContents, &FileContent{
			Content: content,
			Name:    theFile,
		})
	}

	glog.V(4).Infof("Totally %d files read", len(files))

	return
}
