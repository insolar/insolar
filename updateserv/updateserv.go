package updateserv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updater/request"
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

func (updServer *UpdateServer) versionHandler(ver *request.Version) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, err := json.Marshal(ver)
		if err != nil {
			returnError(w, "MARSHAL_ERROR")
		} else {
			fmt.Fprintf(w, string(response))
		}
	})
}

func (updServer *UpdateServer) LoadVersions() *request.Version {

	if updServer.LatestVersion != "" {
		return request.NewVersion(updServer.LatestVersion)
	}
	files, err := ioutil.ReadDir(updServer.UploadPath)
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

func (us *UpdateServer) Start() (err error) {
	ver := us.LoadVersions()
	if ver != nil {
		http.HandleFunc("/latest", us.versionHandler(ver))
		handler := http.FileServer(http.Dir(path.Join(us.UploadPath, ver.Value)))
		http.Handle("/"+ver.Value+"/", http.StripPrefix("/"+ver.Value, handler))
	}
	log.Info("Server started on localhost:" + us.Port + ", use /upload for uploading files and /{version}/{fileName} for downloading files.")
	go func() {
		if err = us.server.ListenAndServe(); err != nil {
			log.Warn("Update server - ", err)
		}
	}()
	return nil
}

func (us *UpdateServer) Stop() error {
	const timeOut = 5
	log.Infof("Shutting down server gracefully ...(waiting for %d seconds)", timeOut)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	err := us.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop UPDATE server")
	}

	return nil
}
