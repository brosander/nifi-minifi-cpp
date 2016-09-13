package properties

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	//"time"

	"github.com/apache/gominifi/api"
)

var nanos = strings.Join([]string{"ns", "nano", "nanos", "nanoseconds"}, "|")
var millis = strings.Join([]string{"ms", "milli", "millis", "milliseconds"}, "|")
var seconds = strings.Join([]string{"s", "sec", "secs", "second", "seconds"}, "|")
var minutes = strings.Join([]string{"m", "min", "mins", "minute", "minutes"}, "|")
var hours = strings.Join([]string{"h", "hr", "hrs", "hour", "hours"}, "|")
var days = strings.Join([]string{"d", "day", "days"}, "|")
var weeks = strings.Join([]string{"w", "wk", "wks", "week", "weeks"}, "|")

var validTimeUnits = strings.Join([]string{nanos, millis, seconds, minutes, hours, days, weeks}, "|")
var validTimeRegex = regexp.MustCompile(`(\d+)\s*(` + validTimeUnits + ")")

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

/*func (p *PropertyValueImpl) AsTimePeriod() time.Duration {
	return p.val
}*/

// Verify PropertyValue is implemented
var _ api.PropertyValue = (*PropertyValueImpl)(nil)
