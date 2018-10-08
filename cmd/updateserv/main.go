/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updater/request"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
)

//const maxUploadSize = 50 * 1024 // 50 MB
const uploadPath = "./data"

func main() {

	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	initLogger(cfgHolder.Configuration.Log)
	port := pflag.StringP("port", "p", "2345", "port to listen")

	//http.HandleFunc("/latest", uploadFileHandler())

	ver := getLatestVersion()
	if ver != "" {
		http.HandleFunc("/latest", versionHandler(ver))
		fs := http.FileServer(http.Dir(path.Join(uploadPath, ver)))
		http.Handle("/"+ver+"/", http.StripPrefix("/"+ver, fs))
	}
	log.Info("Server started on localhost:" + *port + ", use /upload for uploading files and /{version}/{fileName} for downloading files.")
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func getLatestVersion() string {
	files, err := ioutil.ReadDir(path.Join(uploadPath, "latest"))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		return f.Name()
	}
	return ""
}

func versionHandler(ver string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		latest := request.NewVersion(ver)
		response, err := json.Marshal(latest)
		if err != nil {
			returnError(w, "MARSHAL_ERROR")
		} else {
			fmt.Fprintf(w, string(response))
		}
	})
}

// ToDo: create uploader
//func uploadFileHandler() http.HandlerFunc {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// validate file size
//		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
//		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
//			returnError(w, "FILE_TOO_BIG")
//			return
//		}
//
//		// parse and validate file and post parameters
//		fileType := r.PostFormValue("type")
//		file, _, err := r.FormFile("uploadFile")
//		if err != nil {
//			returnError(w, "INVALID_FILE")
//			return
//		}
//		defer file.Close()
//		fileBytes, err := ioutil.ReadAll(file)
//		if err != nil {
//			returnError(w, "INVALID_FILE")
//			return
//		}
//
//		// check file type, detectcontenttype only needs the first 512 bytes
//		filetype := http.DetectContentType(fileBytes)
//		if (filetype != "image/data") {
//			returnError(w, "INVALID_FILE_TYPE")
//			return
//		}
//		fileName := "123"
//		fileEndings, err := mime.ExtensionsByType(fileType)
//		if err != nil {
//			returnError(w, "CANT_READ_FILE_TYPE")
//			return
//		}
//		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
//		fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)
//
//		// write file
//		newFile, err := os.Create(newPath)
//		if err != nil {
//			returnError(w, "CANT_WRITE_FILE")
//			return
//		}
//		defer newFile.Close() // idempotent, okay to call twice
//		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
//			returnError(w, "CANT_WRITE_FILE")
//			return
//		}
//		w.Write([]byte("SUCCESS"))
//	})
//}

func returnError(w http.ResponseWriter, message string) {
	log.Error(message)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(cfg.Level)
	if err != nil {
		log.Errorln(err.Error())
	}
}
