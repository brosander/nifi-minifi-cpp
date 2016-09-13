package properties

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/apache/gominifi/api"
)

func ReadProps(r io.Reader) map[string]string {
	var result = make(map[string]string)
	scanner := bufio.NewScanner(r)
	var buffer bytes.Buffer
	for scanner.Scan() {
		text := scanner.Text()
		buffer.WriteString(text)
		if strings.HasSuffix(text, `\`) {
			buffer.WriteString("\n")
		} else {
			propLine := buffer.String()
			if !strings.HasPrefix(propLine, "#") {
				split := strings.Split(propLine, "=")
				if len(split) == 2 {
					result[split[0]] = split[1]
				}
			}
			buffer.Reset()
		}
	}
	return result
}

func ReadPropsFile(fileName string) (map[string]string, error) {
	propsFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer propsFile.Close()

	return ReadProps(propsFile), nil
}

type PropertyValueImpl struct {
	val string
}

func NewPropertyValue(val string) *PropertyValueImpl {
	return &PropertyValueImpl{val: val}
}

func (p *PropertyValueImpl) AsString() string {
	return p.val
}

func (p *PropertyValueImpl) AsBool() bool {
	result, err := strconv.ParseBool(p.val)
	if err != nil {
		return false
	}
	return result
}

// Verify PropertyValue is implemented
var _ api.PropertyValue = (*PropertyValueImpl)(nil)
