package request

import (
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
)

// Just to make Goland happy
func TestGetProtocol(t *testing.T) {
	assert.Equal(t, getProtocolFromAddress("http://localhost:7087/"), "http", "Get protocol utility success")
	assert.Equal(t, getProtocolFromAddress("localhost:7087"), "", "Get protocol utility success")
}

func TestCompare(t *testing.T) {
	assert.Equal(t, compare(1, 0), 1)
	assert.Equal(t, compare(1, 1), 0)
	assert.Equal(t, compare(0, 1), -1)
}

func TestExtractValue(t *testing.T) {
	value := strings.Split("1.2.3", ".")
	assert.Equal(t, extractIntValue(value, 0), 1)
	assert.Equal(t, extractIntValue(value, 1), 2)
	assert.Equal(t, extractIntValue(value, 2), 3)
}

func TestNewVersion(t *testing.T) {
	v1 := NewVersion("v1.2.3")
	assert.Equal(t, v1.Revision, 3)
	assert.Equal(t, v1.Major, 1)
	assert.Equal(t, v1.Minor, 2)
}

func TestExtractVersion(t *testing.T) {
	v1 := NewVersion("v1.2.3")
	v2 := ExtractVersion("{\"latest\":\"v1.2.3\",\"major\":1,\"minor\":2,\"revision\":3}")
	assert.Equal(t, v2, v1)
	assert.Equal(t, v2.Revision, 3)
	assert.Equal(t, v2.Major, 1)
	assert.Equal(t, v2.Minor, 2)
}

func TestCompareVersion(t *testing.T) {
	v1 := NewVersion("v1.2.3")
	v2 := NewVersion("v1.2.4")
	assert.Equal(t, CompareVersion(v1, v2), -1)
	assert.Equal(t, CompareVersion(v1, v1), 0)
	assert.Equal(t, CompareVersion(v2, v1), 1)
}

func TestGetMaxVersion(t *testing.T) {
	v1 := NewVersion("v1.2.3")
	v2 := NewVersion("v1.2.4")
	assert.Equal(t, GetMaxVersion(v1, v2), v2)
	assert.Equal(t, GetMaxVersion(v2, v1), v2)
	assert.Equal(t, GetMaxVersion(v1, v1), v1)
}

func TestFailGetMaxVersion(t *testing.T) {
	v1 := NewVersion("")
	v2 := NewVersion("v1.2.4")
	v3 := NewVersion("unset")
	assert.Equal(t, GetMaxVersion(v1, v2), v2)
	assert.Equal(t, GetMaxVersion(v2, v1), v2)
	assert.Equal(t, GetMaxVersion(v2, v3), v2)
	assert.Equal(t, GetMaxVersion(v3, v2), v2)
	assert.Equal(t, GetMaxVersion(v1, v3), v3)
}
