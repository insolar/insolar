// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessagebus_MarshalJSON(t *testing.T) {
	testArgs := genArgsData()

	res, err := testArgs.MarshalJSON()
	assert.NoError(t, err)

	assert.Equal(t, ArgumentsString, string(res))
}

func TestMessagebus_ConvertArgs(t *testing.T) {
	testArgs := genArgsData()
	result := make([]interface{}, 0)

	err := convertArgs(testArgs, &result)
	assert.NoError(t, err)

	assert.Equal(t, uint64(0), result[0])

	innerArray := result[1]

	expInnerArray := make([]interface{}, 0)
	expInnerArray = append(expInnerArray, "88cb82ed-1429-4c2c-a963-d283aafe6730")
	expInnerArray = append(expInnerArray, "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYI"+
		"KoZIzj0DAQcDQgAEm2zH7zcB5XsJ5I+Tlcb1gSWaEQME\njtqQ+IWZ6+8LJ1A9Xz2PMnrdRviY1PsA6uYEAkRP/izER"+
		"RCHSx2NGro2pg==\n-----END PUBLIC KEY-----\n")

	assert.Equal(t, expInnerArray, innerArray)
	assert.Equal(t, int64(-21), result[2])
	assert.Equal(t, int64(-1), result[3])
	assert.Equal(t, []interface{}{"CreateMember"}, result[4])
}

func genArgsData() Arguments {
	str := "hVhAAAEAAVl3NEo5/wAJaSpdf7KeSZRV9GSgVHU6ijbZ238AAQABWXc0Sjn/AAlpKl1/sp5JlFX0ZKBUdTqKNtnbf2xDcmVhdGVNZW1i" +
		"ZXJY24J4JDg4Y2I4MmVkLTE0MjktNGMyYy1hOTYzLWQyODNhYWZlNjczMHiyLS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUh" +
		"Lb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFbTJ6SDd6Y0I1WHNKNUkrVGxjYjFnU1dhRVFNRQpqdHFRK0lXWjYrOExKMUE5WHoyUE1ucm" +
		"RSdmlZMVBzQTZ1WUVBa1JQL2l6RVJSQ0hTeDJOR3JvMnBnPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tClggNPKqzYeVtUIvfohh85ppF" +
		"vy7xzwA5ybLGwRZOKBlKtFYQiAgUGHdrUMnhaMPJgIlPgc771kdZl8za8jYfO2rlXAItEm1x6Wv8vbOkD0efySyToEaxvqySnZwppulAWf/" +
		"PJZN0A=="

	decoded, _ := base64.StdEncoding.DecodeString(str)

	return Arguments(decoded)
}

const (
	ArgumentsString = `[0,["88cb82ed-1429-4c2c-a963-d283aafe6730","-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIK` +
		`oZIzj0DAQcDQgAEm2zH7zcB5XsJ5I+Tlcb1gSWaEQME\njtqQ+IWZ6+8LJ1A9Xz2PMnrdRviY1PsA6uYEAkRP/izERRCHSx2NGro2pg==\n` +
		`-----END PUBLIC KEY-----\n"],-21,-1,["CreateMember"]]`
)
