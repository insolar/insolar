package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dgraph-io/badger"

	"github.com/gorilla/mux"
	"github.com/insolar/insolar/log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type MergeJsonRequest struct {
	BkpName string `json:"bkpName"`
	RunGC   bool   `json:"rubGC"`
}

type MergeJsonResponse struct {
	Message string `json:"message"`
}

var bdb *badger.DB

func sendHttpResponse(w http.ResponseWriter, statusCode int, resp MergeJsonResponse) {
	h := w.Header()
	h.Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("sendHttpResonse, json.Marshal: %v\n", err)
		return
	}

	log.Infof("sendHttpResonse: statusCode = %d, resp = %s", statusCode, respBytes)

	_, err = w.Write(respBytes)
	if err != nil {
		log.Errorf("sendHttpResonse, w.Write: %v\n", err)
	}
}

func MergeHttpHandler(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("MergeHttpHandler, ioutil.ReadAll: %v\n", err)
		return
	}

	log.Infof("Processing request: %s", reqBytes)

	var req MergeJsonRequest
	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		log.Errorf("MergeHttpHandler, json.Unmarshal: %v\n", err)
		return
	}

	if req.BkpName == "" {
		sendHttpResponse(w, 400, MergeJsonResponse{
			Message: "Missing bkpName",
		})
		return
	}

	log.Infof("Merging incremental backup, bkpName = %s", req.BkpName)

	bkpFile, err := os.Open(req.BkpName)
	if err != nil {
		sendHttpResponse(w, 400, MergeJsonResponse{
			Message: "Failed to open incremental backup file",
		})
		return
	}
	defer bkpFile.Close()

	log.Info("Backup file is opened")

	err = bdb.Load(bkpFile, 1)
	if err != nil {
		sendHttpResponse(w, 400, MergeJsonResponse{
			Message: "Failed to load incremental backup file",
		})
		return
	}

	msg := "Merge done."
	if req.RunGC {
		err = bdb.RunValueLogGC(0.7)
		if err == nil {
			msg += " GC done."
		} else {
			msg += " GC failed: " + err.Error()
		}
	}

	sendHttpResponse(w, 200, MergeJsonResponse{
		Message: msg,
	})
}

func daemon(listenAddr string, targetDBPath string) {
	log.Info("Opening DB...")

	ops := badger.DefaultOptions(targetDBPath)
	ops.Logger = badgerLogger
	bdb, err := badger.Open(ops)
	if err != nil {
		err = errors.Wrap(err, "failed to open DB")
		exitWithError(err)
	}
	log.Info("DB is opened")

	err = isDBEmpty(bdb)
	if err == nil {
		// will exit
		err = errors.New("DB must not be empty")
		closeRawDB(bdb, err)
	}

	log.Info("DB opened and it's not empty. Starting HTTP server...")
	r := mux.NewRouter().
		PathPrefix("/api/v1").
		Path("/merge").
		Subrouter()
	r.Methods("POST").
		HandlerFunc(MergeHttpHandler)
	http.Handle("/", r)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
	log.Info("HTTP server terminated\n")
}

func parseDaemonParams(ctx context.Context) *cobra.Command {
	var (
		listenAddr   string
		targetDBPath string
	)

	var daemonCmd = &cobra.Command{
		Use:   "daemon",
		Short: "run merge daemon",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Starting merge daemon, address = %s, target-db = %s", listenAddr, targetDBPath)
			daemon(listenAddr, targetDBPath)
		},
	}
	mergeFlags := daemonCmd.Flags()
	targetDBFlagName := "target-db"
	mergeFlags.StringVarP(
		&targetDBPath, targetDBFlagName, "t", "", "directory where backup will be roll to (required)")
	mergeFlags.StringVarP(
		&listenAddr, "address", "a", ":8080", "listen address")

	err := cobra.MarkFlagRequired(mergeFlags, targetDBFlagName)
	if err != nil {
		err := errors.Wrap(err, "failed to set required param: "+targetDBFlagName)
		exitWithError(err)
	}

	return daemonCmd
}

func daemonMerge(address string, backupFileName string, runGC bool) {
	reqJson := MergeJsonRequest{BkpName: backupFileName, RunGC: runGC}
	reqBytes, err := json.Marshal(reqJson)
	if err != nil {
		err = errors.Wrap(err, "daemonMerge - json.Marshal failed")
		exitWithError(err)
	}

	req, err := http.NewRequest("POST", address+"/api/v1/merge", bytes.NewBuffer(reqBytes))
	if err != nil {
		err = errors.Wrap(err, "daemonMerge - http.NewRequest failed")
		exitWithError(err)
	}

	client := http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "daemonMerge - client.Do failed")
		exitWithError(err)
	}
	defer httpResp.Body.Close()

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		err = errors.Wrap(err, "daemonMerge - ioutil.ReadAll failed")
		exitWithError(err)
	}
	var resp MergeJsonResponse
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		err = errors.Wrap(err, "daemonMerge - json.Unmarshal failed")
		exitWithError(err)
	}

	if httpResp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Merge failed: daemon returned code %d and body: %s\n", httpResp.StatusCode, respBytes))
		exitWithError(err)
	}

	log.Infof("HTTP response OK. Daemon: %s", resp.Message)
}

func parseDaemonMergeParams(ctx context.Context) *cobra.Command {
	var (
		address        string
		backupFileName string
		runGC          bool
	)

	var daemonMergeCmd = &cobra.Command{
		Use:   "daemon-merge",
		Short: "merge incremental backup using merge daemon",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Starting daemon-merge, address = %s, bkp-name = %s", address, backupFileName)
			daemonMerge(address, backupFileName, runGC)
		},
	}
	mergeFlags := daemonMergeCmd.Flags()
	bkpFileName := "bkp-name"
	mergeFlags.StringVarP(
		&backupFileName, bkpFileName, "n", "", "file name if incremental backup (required)")
	mergeFlags.StringVarP(
		&address, "address", "a", "http://localhost:8080", "merge daemon listen address")
	mergeFlags.BoolVarP(
		&runGC, "run-gc", "g", false, "run GC")

	err := cobra.MarkFlagRequired(mergeFlags, bkpFileName)
	if err != nil {
		err := errors.Wrap(err, "failed to set required param: "+bkpFileName)
		exitWithError(err)
	}

	return daemonMergeCmd
}
