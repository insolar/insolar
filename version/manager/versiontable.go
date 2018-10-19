/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the License);
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an AS IS BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */
package manager

const (
	INS_BASE_VERSION  = "v0.3.0"
	INS_VERSION_TABLE = "#  Structure of the versiontable.yml file" +
		"#" +
		"#versiontable:" +
		"#  {FEATURE KEYWORD}:" +
		"#    startversion: {Start version}" +
		"#    description: {Description about feature}" +
		"#" +
		"#" +
		"#For example, the file structure can be:" +
		"#" +
		"#versiontable:" +
		"#  CREATED_VERSION_MANAGER:" +
		"#    startversion: v0.1.1" +
		"#    description: Create version manager for Insolar platform" +
		"#  CHANGED_CONSENSUS_VERSION_MANAGER:" +
		"#    startversion: v0.5.2" +
		"#    description: Changed consensus of the version manager to BFT" +
		"#..." +
		"" +
		"versiontable:" +
		""
)
