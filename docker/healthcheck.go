package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetEnvDefault(key, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}

func obtainContainerRole() string {
	return strings.ToLower(GetEnvDefault("INSOLARD_ROLE", "insolard+insgorund"))
}

func obtainDockerPublicIP() string {
	cmd := exec.Command("awk", "END{print $1}", "/etc/hosts")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Print("Failed to obtain public IP:", err.Error())
		fmt.Print(stderr.String())
		os.Exit(1)
	}

	return strings.Trim(stdout.String(), " \n\t")
}

func getURI(port uint) string {
	return fmt.Sprintf("%s:%d", obtainDockerPublicIP(), port)
}

const (
	defaultApiListenPort = 19101
)

const statusBody = "{\"jsonrpc\": \"2.0\", \"method\": \"status.Get\", \"id\": 0}"

func checkInsolard() int {
	var apiUrl url.URL
	apiUrl.Scheme = "http"
	apiUrl.Host = GetEnvDefault("INSOLARD_API_LISTEN", getURI(defaultApiListenPort))
	apiUrl.Path = "/api/rpc"

	client := http.Client{
		Transport: &http.Transport{},
		Timeout:   60 * time.Second,
	}

	body := strings.NewReader(statusBody)
	res, err := client.Post(apiUrl.String(), "application/json", body)
	if err != nil {
		fmt.Printf("Failed to make HTTP request:", err.Error())
		return 1
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read body:", err.Error())
		return 1
	}
	var out struct {
		Result struct {
			PulseNumber    uint32
			NetworkState   string
			NodeState      string
			Origin         struct{ Role string }
			ActiveListSize int
		}
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		fmt.Printf("Failed to parse body:", err.Error())
		return 1
	}
	// TODO: what to check in output ?
	fmt.Print(data)
	return 0
}

func checkInsgorund() int {
	// TODO: implement healthcheck contract calling
	return 0
}

func main() {
	role := obtainContainerRole()

	retcode := 0
	if strings.Index(role, "insolard") != -1 {
		retcode |= checkInsolard()
	}
	if strings.Index(role, "insgorund") != -1 {
		retcode |= checkInsgorund()
	}
	os.Exit(retcode)
}
