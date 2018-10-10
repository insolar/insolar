package updater

import (
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updater/request"
	"github.com/insolar/insolar/version"
)

type Updater struct {
	ServersList       []string
	BinariesList      []string
	LastSuccessServer string
	CurrentVer        string
	Delay             int64
}

func NewUpdater() *Updater {
	return &Updater{
		[]string{"http://localhost:2345"},
		[]string{"insgocc", "insgorund", "insolar", "insolard", "pulsard", "insupdater"},
		"",
		version.Version,
		60,
	}
}

func (updater *Updater) IsSameVersion(currentVersion string) (bool, string, error) {
	log.Debug("Verify latest peer version from remote server")
	updater.CurrentVer = currentVersion
	currentVer := request.NewVersion(currentVersion)
	if updater.LastSuccessServer != "" {
		log.Debug("Latest update server was: ", updater.LastSuccessServer)
		vers, err := request.ReqCurrentVerFromAddress(request.GetProtocol(updater.LastSuccessServer), updater.LastSuccessServer)
		if err == nil && vers != "" {
			versionFromUS := request.ExtractVersion(vers)
			return request.CompareVersion(versionFromUS, currentVer) < 0, versionFromUS.Value, nil
		}
	}
	lastSuccessServer, versionFromUS, err := request.ReqCurrentVer(updater.ServersList)
	if err != nil {
		return true, "", err
	}
	log.Debug("Get version=", versionFromUS.Value, " from remote server: ", lastSuccessServer)
	updater.LastSuccessServer = lastSuccessServer

	if versionFromUS == nil || updater.CurrentVer == "" {
		return true, "unset", nil
	} else
	//if(updater.currentVer != versionFromUS){
	if request.CompareVersion(versionFromUS, currentVer) > 0 {
		return false, versionFromUS.Value, nil
	}
	return true, versionFromUS.Value, nil
}

func (updater Updater) DownloadFiles(version string) (success bool) {
	log.Info("Start download files from remote server")
	if updater.LastSuccessServer == "" {
		return false
	}
	return request.DownloadFiles(version, updater.BinariesList, updater.LastSuccessServer)
}
