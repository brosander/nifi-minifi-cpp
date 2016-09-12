package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("getfile")

type GetFile struct {
	recursive      bool
	minSize        int64
	maxSize        int64
	minAge         time.Time
	maxAge         time.Time
	ignoreHidden   bool
	keepSourceFile bool
}

func (getFile *GetFile) OnTrigger(context *ProcessContext, session *ProcessSession) error {
	files, err := getFile.performListing("/tmp/mgo")
	if err != nil {
		return err
	}
	for fileName := range files {
		file := getFile.open(fileName)
		if file == nil {
			continue
		}
		deleteName := fileName
		defer func() {
			file.Close()
			if !getFile.keepSourceFile {
				err := os.Remove(deleteName)
				if err != nil {
					log.Error("Unable to delete " + deleteName + " even though we could open it read/write. (" + err.Error() + ")")
				}
			}
		}()
	}
	return nil
}

func (getFile *GetFile) open(fileName string) *os.File {
	if getFile.keepSourceFile {
		file, err := os.Open(fileName)
		if err != nil {
			log.Warning("Unable to open file for reading: " + fileName + " (" + err.Error() + ")")
			return nil
		}
		return file
	} else {
		file, err := os.OpenFile(fileName, os.O_RDWR, 0600)
		if err != nil {
			log.Warning("Unable to open file for read/write: " + fileName + " (" + err.Error() + ")")
			return nil
		}
		return file
	}
}

func (getFile *GetFile) performListing(dir string) (chan string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	files := make(chan string, 100)
	go func() {
		getFile.doPerformListing(absDir, files)
		close(files)
	}()
	return files, nil
}

func (getFile *GetFile) doPerformListing(absoluteDirectory string, files chan string) {
	dirFiles, err := ioutil.ReadDir(absoluteDirectory)
	if err != nil {
		log.Warning("Unable to read directory: " + absoluteDirectory + " (" + err.Error() + ")")
		return
	}
	for _, file := range dirFiles {
		if file.IsDir() {
			if getFile.recursive {
				getFile.doPerformListing(filepath.Join(absoluteDirectory, file.Name()), files)
			}
		} else {
			absPath := filepath.Join(absoluteDirectory, file.Name())
			if getFile.acceptFile(absPath, file) {
				files <- absPath
			}
		}
	}
}

func (getFile *GetFile) acceptFile(absolutePath string, file os.FileInfo) bool {
	if getFile.minSize > 0 && file.Size() < getFile.minSize {
		return false
	}
	if getFile.maxSize > 0 && file.Size() > getFile.maxSize {
		return false
	}
	if !getFile.minAge.IsZero() && file.ModTime().Before(getFile.minAge) {
		return false
	}
	if !getFile.maxAge.IsZero() && file.ModTime().After(getFile.maxAge) {
		return false
	}
	return true
}
