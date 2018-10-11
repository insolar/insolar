package updater

import (
	"os"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updater/request"
	"github.com/insolar/insolar/version"
)

func (up *Updater) verifyAndUpdate() error {
	log.Info("Try verify for update ")
	sameVersion, newVersion, err := up.IsSameVersion(version.Version)
	if err != nil {
		return err
	}
	if !sameVersion {
		log.Debug("Current version: ", version.Version, ", found version: ", newVersion)
		// Run Update
		if up.DownloadFiles(newVersion) {
			// ToDo: send stop signal, then copy files from folder=./${VERSION} to current folder

			os.Setenv("INS_LATEST_VER", newVersion)
		}
	}
	// Run peer
	//executePeer()
	// ToDo: Run update service with timer
	// exit
	return nil
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

func (up *Updater) DownloadFiles(version string) (success bool) {
	if up.started {
		return false
	}
	log.Info("Start download files from remote server")
	if up.LastSuccessServer == "" {
		return false
	}
	up.started = true
	success = request.DownloadFiles(version, up.BinariesList, up.LastSuccessServer)
	up.started = false
	return
}
