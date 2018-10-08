package request

import (
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type HttpRequestUpdateNode struct {
	RequestUpdateNode
}

func (request HttpRequestUpdateNode) getCurrentVer(address string) (string, error) {
	response, err := http.Get(address + "/latest")
	if err != nil {
		return "Error during http request", err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "Error, cannot get response body", err
	}
	return string(content), nil
}

func (request HttpRequestUpdateNode) downloadFile(filepath string, url string) error {

	//Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Error("OS Create file error: ", err)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Error("HTTP server error: ", err)
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Error("HTTP bad status: ", resp.Status)
		return errors.Errorf("HTTP error, ", resp.Status)
	}

	// Writer the body to file
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Error("OS write file error: ", err)
		return err
	}
	log.Info("Downloaded file: "+url+", save to: "+filepath+", total bytes: ", written)
	return nil
}
