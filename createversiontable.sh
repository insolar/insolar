#!/bin/bash

FILE=$1
GO_FILE=$2
INS_BASE_VERSION=v0.3.0

echo -e "/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */" > $GO_FILE

echo -e "package manager

const (
	INS_BASE_VERSION = \"$INS_BASE_VERSION\"
	INS_VERSION_TABLE = " >> $GO_FILE


while read LINE; do
     echo -e "\"$LINE\"+" >> $GO_FILE
done < $FILE
echo -e "\"\")" >> $GO_FILE

