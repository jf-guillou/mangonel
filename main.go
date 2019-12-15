package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
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
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(configuration.StoragePath)))

	log.Info("Listen and serve")
	err := http.ListenAndServe(configuration.ListenAddr, router)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Bye")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var filename string
	var err error

	reader, err := r.MultipartReader()
	if err != nil {
		log.Error("Unable to load multipart reader : ", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	log.Debug("Reading multipart data")
	for {
		part, err := reader.NextPart()
		if err == io.EOF || part == nil {
			break
		}

		switch part.FormName() {
		case "filepond":
			if part.FileName() == "" {
				continue
			}

			err = checkFileSize(r.Header)
			if err != nil {
				log.Error("File size checks failed : ", err)
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			filename, err = handleFilePart(part)
			if err != nil {
				log.Error("File read failed : ", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}
	}

	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(filename))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/index.html")
}
