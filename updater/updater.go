package updater

import (
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updater/request"
	"github.com/insolar/insolar/version"
)

type Updater struct {
	serversList       []string
	binariesList      []string
	lastSuccessServer string
	currentVer        string
	Delay             int64
}

func NewUpdater() *Updater {
	newUpdater := Updater{}
	newUpdater.serversList = []string{"http://localhost:2345"}
	newUpdater.lastSuccessServer = ""
	newUpdater.binariesList = []string{"insgocc", "insgorund", "insolar", "insolard", "pulsard", "insupdater"}
	newUpdater.currentVer = version.Version
	newUpdater.Delay = 60
	log.Debug("Create new Updater: ", newUpdater)
	return &newUpdater
}

func (updater *Updater) IsSameVersion(currentVersion string) (bool, string, error) {
	log.Debug("Verify latest peer version from remote server")
	updater.currentVer = currentVersion
	currentVer := request.NewVersion(currentVersion)
	if updater.lastSuccessServer != "" {
		log.Debug("Latest update server was: ", updater.lastSuccessServer)
		version, err := request.ReqCurrentVerFromAddress(request.GetProtocol(updater.lastSuccessServer), updater.lastSuccessServer)
		if err == nil && version != "" {
			versionFromUS := request.ExtractVersion(version)
			return request.CompareVersion(versionFromUS, currentVer) < 0, versionFromUS.Value, nil
		}
	}
	lastSuccessServer, versionFromUS, err := request.ReqCurrentVer(updater.serversList)
	log.Debug("Get version=", versionFromUS.Value, " from remote server: ", lastSuccessServer)
	updater.lastSuccessServer = lastSuccessServer
	if err != nil {
		return true, versionFromUS.Value, err
	}
	if versionFromUS == nil || updater.currentVer == "" {
		return true, "unset", nil
	} else
	//if(updater.currentVer != versionFromUS){
	if request.CompareVersion(versionFromUS, currentVer) > 0 {
		return false, versionFromUS.Value, nil
	}
	return true, versionFromUS.Value, nil
}

func (updater Updater) DownloadFiles(version string) bool {
	log.Info("Start download files from remote server")
	if updater.lastSuccessServer != "" {
		return request.DownloadFiles(version, updater.binariesList, updater.lastSuccessServer)
	} else {
		return false
	}

}
