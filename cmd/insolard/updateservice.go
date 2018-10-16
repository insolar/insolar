package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/insolar/insolar/cmd/insolard/componentmanager"
	"github.com/insolar/insolar/updater"
)

func startUpdateProcess(newVersion string) {
	componentmanager.GetComponentManager().StopAll()
	cmd := exec.Command(path.Join(".", "bin", newVersion, "insolard"), "--config", path.Join(".", "scripts", "insolard", "insolard.yaml"))
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if out != nil {
		fmt.Println(out)
	}
	err = addValueToEnvFile("INS_LATEST_VER", newVersion)
	if err != nil {
		fmt.Println("Can not set OS envelop value INS_LATEST_VER (.env): ", err)
	}
	err = os.Setenv("INS_LATEST_VER", newVersion)
	if err != nil {
		fmt.Println("Can not set OS envelop value INS_LATEST_VER: ", err)
	}
	os.Exit(0)
}

func runUpdateService(updater *updater.Updater) {
	fmt.Println("Starting UPDATE service")
	for {
		time.Sleep(time.Minute)
		if updater.ReadyToUpdate {
			fmt.Println("Update service is ready to restart insolard")
			startUpdateProcess(updater.CurrentVer)
		}
	}
}

func addValueToEnvFile(key string, value string) error {
	path := path.Join(".", "scripts", "insolard", ".env")
	envs, err := readLines(path)
	if err != nil {
		fmt.Println(err.Error())
	}

	for index, env := range envs {
		if strings.Contains(env, key+"=") {
			envs[index] = key + "=" + value
			return writeLines(envs, path)
		}
	}
	envs = append(envs, key+"="+value)
	return writeLines(envs, path)

	// fileHandle, _ := os.OpenFile(path.Join(".","scripts","insolard",".env"), os.O_RDWR, 0666)
	// writer := bufio.NewWriter(fileHandle)
	// defer fileHandle.Close()
	//
	// fmt.Fprintln(writer, key+"="+value)
	// writer.Flush()

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

// INS_LATEST_VER=
// VER=$([ "$INS_LATEST_VER" ] && echo "$INS_LATEST_VER" || echo $(git describe --abbrev=0 --tags))
// BIN_DIR=bin/$VER
// TEST_DATA=testdata
// INSOLARD=$BIN_DIR/insolard
// INSGORUND=$BIN_DIR/insgorund
// CONTRACT_STORAGE=contractstorage
// LEDGER_DIR=data
// INSGORUND_LISTEN_PORT=18181
// INSGORUND_RPS_PORT=18182
