package request

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

// Just to make Goland happy
func TestStubCurrentVer(t *testing.T) {
	newVer := NewVersion("v1.2.3")
	assert.Equal(t, newVer.Major, 1, "Major verify passed")
	assert.Equal(t, newVer.Minor, 2, "Minor verify passed")
	assert.Equal(t, newVer.Revision, 3, "Revision verify passed")
}
