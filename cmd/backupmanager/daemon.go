package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insolar/insolar/log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type MergeJsonRequest struct {
	BkpName string `json:"bkpName"`
}

type MergeJsonResponse struct {
	Message string `json:"message"`
}

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

	// AALEKSEEV TODO actually process req.BkpName

	sendHttpResponse(w, 200, MergeJsonResponse{
		Message: "Merge done",
	})
}

func daemon(listenAddr string, targetDBPath string) {
	r := mux.NewRouter().
		PathPrefix("/api/v1").
		Path("/merge").
		Subrouter()
	r.Methods("POST").
		HandlerFunc(MergeHttpHandler)
	http.Handle("/", r)
	err := http.ListenAndServe(listenAddr, nil)
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

func parseDaemonMergeParams(ctx context.Context) *cobra.Command {
	var (
		daemonHost     string
		daemonPort     int
		backupFileName string
	)

	var daemonMergeCmd = &cobra.Command{
		Use:   "daemon-merge",
		Short: "merge incremental backup using merge daemon",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Starting daemon-merge, host = %s, port = %d, bkp-name = %s", daemonHost, daemonPort, backupFileName)
			// daemonMerge(daemonHost, daemonPort, backupFileName) // AALEKSEEV TODO
		},
	}
	mergeFlags := daemonMergeCmd.Flags()
	bkpFileName := "bkp-name"
	mergeFlags.StringVarP(
		&backupFileName, bkpFileName, "n", "", "file name if incremental backup (required)")
	mergeFlags.StringVarP(
		&daemonHost, "address", "a", "localhost", "merge daemon listen address or host")
	mergeFlags.IntVarP(
		&daemonPort, "port", "p", 8080, "merge daemon listen port")

	err := cobra.MarkFlagRequired(mergeFlags, bkpFileName)
	if err != nil {
		err := errors.Wrap(err, "failed to set required param: "+bkpFileName)
		exitWithError(err)
	}

	return daemonMergeCmd
}
