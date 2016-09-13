package getfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/apache/gominifi/api"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("getfile")

var success = api.NewRelationship("success", "All files are routed to success")

var supportedProperties = []api.Property{}
var supportedRelationships = []api.Relationship{success}

type GetFile struct {
	name           string
	id             string
	isRunning      bool
	directory      string
	recursive      bool
	minSize        int64
	maxSize        int64
	minAge         time.Time
	maxAge         time.Time
	ignoreHidden   bool
	keepSourceFile bool

	readDir func(string) ([]os.FileInfo, error)
}

func NewGetFile(name string, id string) *GetFile {
	return &GetFile{
		name:    name,
		id:      id,
		readDir: ioutil.ReadDir,
	}
}

func (getFile *GetFile) Name() string {
	return getFile.name
}

func (getFile *GetFile) Id() string {
	return getFile.id
}

func (getFile *GetFile) IsRunning() bool {
	return getFile.isRunning
}

func (getFile *GetFile) SupportedProperties() []api.Property {
	return supportedProperties
}

func (getFile *GetFile) SupportedRelationships() []api.Relationship {
	return supportedRelationships
}

func (getFile *GetFile) OnTrigger(context api.ProcessContext, session api.ProcessSession) error {
	files, err := getFile.performListing("/tmp/mgo")
	if err != nil {
		return err
	}
	for fileName := range files {
		attributes := make(map[string]string)
		attributes["filename"] = filepath.Base(fileName)
		attributes["path"] = filepath.Dir(fileName)
		attributes["absolute.path"] = fileName
		var flowFile = session.Create()
		flowFile = session.PutAllAttributes(flowFile, attributes)
		flowFile = session.ImportFrom(fileName, getFile.keepSourceFile, flowFile)
		session.Transfer(flowFile, success)
	}
	return nil
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
	dirFiles, err := getFile.readDir(absoluteDirectory)
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

// Verify that Processor interface is implemented
var _ api.Processor = (*GetFile)(nil)
