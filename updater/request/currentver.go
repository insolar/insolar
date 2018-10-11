package request

import (
	"errors"
	"regexp"

	"github.com/insolar/insolar/log"
)

type Version struct {
	Value    string `json:"latest"`
	Major    int    `json:"major"`
	Minor    int    `json:"minor"`
	Revision int    `json:"revision"`
}

func NewVersion(ver string) *Version {
	v := Version{}
	v.Value = ver
	re := regexp.MustCompile("[0-9]+")
	arr := re.FindAllString(ver, -1)
	v.Major = extractIntValue(arr, 0)
	v.Minor = extractIntValue(arr, 1)
	v.Revision = extractIntValue(arr, 2)
	return &v
}

func ReqCurrentVer(addresses []string) (string, *Version, error) {
	log.Debug("Found update server addresses: ", addresses)

	for _, address := range addresses {
		if address != "" {
			log.Info("Found update server address: ", address)
			ver, err := ReqCurrentVerFromAddress(GetProtocol(address), address)

			if err == nil && ver != "" {
				currentVer := ExtractVersion(ver)
				return address, currentVer, err
			}
		}
	}
	log.Warn("No Update Servers available")
	return "", nil, errors.New("No Update Servers available")
}

func ReqCurrentVerFromAddress(request UpdateNode, address string) (string, error) {
	log.Debug("Check latest version info from remote server: ", address)
	if request == nil {
		return "", errors.New("Unknown protocol")
	}
	return request.getCurrentVer(address)
}
