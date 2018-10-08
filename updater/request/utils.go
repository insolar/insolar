package request

import (
	"encoding/json"
	"fmt"
	"github.com/insolar/insolar/log"
	"os"
	"path"
	"strconv"
	"strings"
)

func GetProtocol(address string) RequestUpdateNode {
	protocol := getProtocolFromAddress(address)
	if protocol != "" {
		switch protocol {
		case "http":
			{
				return HttpRequestUpdateNode{}
			}
		case "ftp":
			{
				return FtpRequestUpdateNode{}
			}
		case "smb":
			{
				return SmbRequestUpdateNode{}
			}
			log.Warn("Unknown protocol ", protocol[0])
		}
	}
	return nil
}

func getProtocolFromAddress(address string) string {
	protocol := strings.Split(address, "://")
	if len(protocol) < 2 {
		return ""
	}
	return protocol[0]
}

func createCurrentPath(version string) string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	pathToSave := path.Join(pwd, version)
	if err := os.Mkdir(pathToSave, 0777); err != nil {
		log.Warn("Error while create folder: ", err)
	}
	return pathToSave
}

func ExtractVersion(ver string) *version {
	latestVersion := version{}
	err := json.Unmarshal([]byte(ver), &latestVersion)
	if err != nil {
		log.Warn("Error parsing data: ", err)
		return nil
	}
	return &latestVersion
}

// if ver1 < ver2   return -1
// if ver1 == ver2  return 0
// if ver1 > ver2   return 1
func CompareVersion(ver1 *version, ver2 *version) (result int) {
	result = 0
	if result = compare(ver1.Major, ver2.Major); result == 0 {
		if result = compare(ver1.Minor, ver2.Minor); result == 0 {
			result = compare(ver1.Revision, ver2.Revision)
		}
	}
	return
}

func compare(x int, y int) int {
	if x < y {
		return -1
	} else if x > y {
		return 1
	} else {
		return 0
	}
}

func extractIntValue(arr []string, index int) int {
	if len(arr) >= index+1 && arr[index] != "" {
		value, err := strconv.Atoi(arr[index])
		if err == nil {
			return value
		}
	}
	return 0
}
