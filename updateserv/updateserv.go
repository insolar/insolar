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
package updateserv

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updater/request"
	"github.com/pkg/errors"
)

type UpdateServer struct {
	UploadPath    string
	Port          string
	LatestVersion string
	server        *http.Server
	//const maxUploadSize = 50 * 1024 // 50 MB
}

func NewUpdateServer(port string, upPath string) *UpdateServer {
	return &UpdateServer{
		upPath,
		port,
		"v0.0.0",
		&http.Server{Addr: ":" + port},
	}
}

func (ups *UpdateServer) versionHandler(ver *request.Version) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, err := json.Marshal(ver)
		if err != nil {
			returnError(w, "MARSHAL_ERROR")
		} else {
			fmt.Fprintf(w, string(response))
		}
	})
}

func (ups *UpdateServer) LoadVersions() *request.Version {

	if ups.LatestVersion != "" {
		return request.NewVersion(ups.LatestVersion)
	}
	files, err := ioutil.ReadDir(ups.UploadPath)
	if err != nil {
		log.Error(err)
		return nil
	}
	newVer := request.NewVersion("v0.0.0")
	for _, f := range files {
		if f.IsDir() {
			newVer = request.GetMaxVersion(newVer, request.NewVersion(f.Name()))
		}
	}
	if newVer.Value != "v0.0.0" {
		return newVer
	}
	return nil
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
	_, err := w.Write([]byte(message))
	log.Error(err)
}

func (ups *UpdateServer) Start() (err error) {
	ver := ups.LoadVersions()
	if ver != nil {
		http.HandleFunc("/latest", ups.versionHandler(ver))
		handler := http.FileServer(http.Dir(path.Join(ups.UploadPath, ver.Value)))
		http.Handle("/"+ver.Value+"/", http.StripPrefix("/"+ver.Value, handler))
	}
	log.Info("Server started on localhost:" + ups.Port + ", use /upload for uploading files and /{version}/{fileName} for downloading files.")
	go func() {
		if err = ups.server.ListenAndServe(); err != nil {
			log.Warn("Update server - ", err)
		}
	}()
	return nil
}

func (ups *UpdateServer) Stop() error {
	const timeOut = 5
	log.Infof("Shutting down server gracefully ...(waiting for %d seconds)", timeOut)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	err := ups.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop UPDATE server")
	}

	return nil
}
