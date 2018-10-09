package request

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

// Just to make Goland happy
func TestGetProtocol(t *testing.T) {
	assert.Equal(t, getProtocolFromAddress("http://localhost:7087/"), "http", "Get protocol utility success")
	assert.Equal(t, getProtocolFromAddress("ftp://localhost:7087/"), "ftp", "Get protocol utility success")
	assert.Equal(t, getProtocolFromAddress("localhost:7087"), "", "Get protocol utility success")
}

func TestExtractVersion(t *testing.T) {
	assert.Equal(t, ExtractVersion("{\"latest\":\"v0.3.1\",\"major\":0,\"minor\":3,\"revision\":1}"), NewVersion("v0.3.1"), "ExtractVersion success")

}
