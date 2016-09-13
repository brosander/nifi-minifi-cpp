package properties

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReadProps(t *testing.T) {
	var props = ReadProps(strings.NewReader("#comment\nkey1=value1\nkey2=value2\n\nkey3=value3\nkey4=multiline\\\nvalue4\nkey5=multiline\\\nvalue5\n\nkey6=multiline\\\nvalue6"))
	assert.Equal(t, 6, len(props))
	assert.Equal(t, "value1", props["key1"])
	assert.Equal(t, "value2", props["key2"])
	assert.Equal(t, "value3", props["key3"])
	assert.Equal(t, "multiline\\\nvalue4", props["key4"])
	assert.Equal(t, "multiline\\\nvalue5", props["key5"])
	assert.Equal(t, "multiline\\\nvalue6", props["key6"])
}
