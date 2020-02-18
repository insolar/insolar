//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

// +build requestertest

package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/insolar/insolar/application/cmd/requester/cmd"
	"github.com/insolar/insolar/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type requesterSuiteTest struct {
	suite.Suite
	ts          *httptest.Server
	paramsFile  *os.File
	userKeyFile *os.File
}

func (s *requesterSuiteTest) SetupSuite() {
	s.ts = httptest.NewServer(handlers())
	s.paramsFile = getRequestParamsFile()

	tempFile, err := ioutil.TempFile("", "requester-test-")
	if err != nil {
		log.Fatal("failed open tmp paramsFile:", err)
	}
	userKeyPath := tempFile.Name()
	writePrivateKeyToFile(userKeyPath)
	s.userKeyFile = tempFile
}

func (s *requesterSuiteTest) TearDownSuite() {
	s.ts.Close()

	_ = s.paramsFile.Close()
	_ = os.Remove(s.paramsFile.Name())

	_ = s.userKeyFile.Close()
	_ = os.Remove(s.userKeyFile.Name())
}

func (s *requesterSuiteTest) TestRequester_HelpWorks() {
	sout, _ := runCmd("--help")
	assert.Contains(s.T(), sout, cmd.ApplicationShortDescription)
}

func (s *requesterSuiteTest) TestRequester_AllArgsPassedSuccessfully() {
	sout, err := runCmd(getArgs(s.userKeyFile.Name(), s.ts.URL, s.paramsFile.Name(), true, true)...)
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), sout, "requestReference")
}

func (s *requesterSuiteTest) TestRequester_UserShouldPassUrl() {
	sout, err := runCmd(getArgs(s.userKeyFile.Name(), "", s.paramsFile.Name(), true, true)[1:]...)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), sout, "The program required url as an argument")
}

func (s *requesterSuiteTest) TestRequester_UserShouldNotPassInvalidUrl() {
	sout, err := runCmd(getArgs("/somePath", "http://localhost:-1", "someParams.json", true, true)...)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), sout, "URL parameter is incorrect")
}

func (s *requesterSuiteTest) TestRequester_UserShouldPassMemberKey() {
	sout, err := runCmd(getArgs("/path/to/nonexisting/key", s.ts.URL, "someParams.json", true, true)...)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), sout, "Member keys does not exists")
}

func (s *requesterSuiteTest) TestRequester_UserShouldPassRequestParams() {
	sout, err := runCmd(getArgs(s.userKeyFile.Name(), s.ts.URL, "/path/to/nonexisting/requestparams", true, true)...)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), sout, "Cannot unmarshal request")
}

func (s *requesterSuiteTest) TestRequester_UserCanPassRequestParamsJson() {
	sout, err := runCmd(getArgs(s.userKeyFile.Name(), s.ts.URL, createMemberRequestExample, true, true)...)
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), sout, "requestReference")
}

func TestAllRequesterTests(t *testing.T) {
	suite.Run(t, new(requesterSuiteTest))
}
