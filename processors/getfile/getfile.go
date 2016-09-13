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

var batchSize = &api.Property{"Batch Size", "The maximum number of files to pull in each iteration", "10"}
var inputDirectory = &api.Property{"Input Directory", "The input directory from which to pull files", "."}
var ignoreHidden = &api.Property{"Ignore Hidden Files", "Indicates whether or not hidden files should be ignored", "true"}
var keepSource = &api.Property{"Keep Source File", "If true, the file is not deleted after it has been copied to the Content Repository", "false"}
var maxAge = &api.Property{"Maximum File Age", "The minimum age that a file must be in order to be pulled; any file younger than this amount of time (according to last modification date) will be ignored", "0 sec"}
var minAge = &api.Property{"Minimum File Age", "The maximum age that a file must be in order to be pulled; any file older than this amount of time (according to last modification date) will be ignored", "0 sec"}
var maxSize = &api.Property{"Maximum File Size", "The maximum size that a file can be in order to be pulled", "0 B"}
var minSize = &api.Property{"Minimum File Size", "The minimum size that a file must be in order to be pulled", "0 B"}
var pollInterval = &api.Property{"Polling Interval", "Indicates how long to wait before performing a directory listing", "0 sec"}
var recursive = &api.Property{"Recurse Subdirectories", "Indicates whether or not to pull files from subdirectories", "true"}

var supportedProperties = []*api.Property{batchSize, inputDirectory, ignoreHidden, keepSource, maxAge, minAge, maxSize, minSize, pollInterval, recursive}

var success = &api.Relationship{"success", "All files are routed to success"}

var supportedRelationships = []*api.Relationship{success}

type GetFile struct {
	api.BaseProcessor

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
	result := &GetFile{
		BaseProcessor: *api.NewBaseProcessor(
			name,
			id,
			supportedProperties,
			supportedRelationships,
		),
		readDir: ioutil.ReadDir,
	}
	return result
}

func (getFile *GetFile) OnTrigger(context api.ProcessContext, session api.ProcessSession) error {
	getFile.recursive = context.GetPropertyValue(recursive).AsBool()
	files, err := getFile.performListing(context.GetPropertyValue(inputDirectory).AsString())
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
