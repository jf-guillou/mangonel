package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var allowedFileExtensions = []string{
	".jpg",
	".jpeg",
	".gif",
	".png",
}

func main() {
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel

	loadConfiguration()
	log.Info("Loading routes")

	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/up", uploadHandler).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/static"))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(configuration.Storage)))

	log.Info("Listen and serve")
	err := http.ListenAndServe(configuration.Addr, router)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Bye")
}

type response struct {
	Success      bool     `json:"success"`
	Error        string   `json:"error"`
	PreventRetry bool     `json:"preventRetry"`
	Reset        bool     `json:"reset"`
	Filenames    []string `json:"filenames"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	files := make([]string, 0)
	validFileSize := false

	reader, err := r.MultipartReader()
	if err != nil {
		jsonErr(w, err)
		return
	}

	log.Debug("Reading multipart data")
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		switch part.FormName() {
		// case "qquid":
		// case "qqfilename":
		case "qqtotalfilesize":
			b, err := ioutil.ReadAll(part)
			if err != nil {
				jsonErr(w, err)
				return
			}

			if configuration.MaxFilesize > 0 {
				filesize, err := strconv.Atoi(string(b))
				if err != nil {
					jsonErr(w, errors.New("malformed filesize"))
					return
				}

				if filesize > configuration.MaxFilesize {
					jsonErr(w, errors.New("file is too big (>"+string(configuration.MaxFilesize)+")"))
					return
				}
			}
			validFileSize = true
		case "qqfile":
			if !validFileSize {
				jsonErr(w, errors.New("unable to determine file size"))
				return
			}

			filename, err := handleFilePart(part)
			if err != nil {
				jsonErr(w, err)
				return
			}

			files = append(files, filename)
			validFileSize = false
		}
	}

	err = json.NewEncoder(w).Encode(response{
		Success:   true,
		Filenames: files,
	})

	if err != nil {
		log.Error("Encode errored", err)
	}
}

func jsonErr(w http.ResponseWriter, e error) {
	log.Warn("Request errored ", e)
	err := json.NewEncoder(w).Encode(response{
		Success: false,
		Error:   e.Error(),
	})

	if err != nil {
		log.Error("Encode errored ", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/index.html")
}
