package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
)

// Configuration read from config.json
type Configuration struct {
	// Hash length
	Length int
	// Listening address:port
	Addr string
	// Storage path
	Storage string
}

var configuration Configuration

const minHashLen int = 1
const maxHashLen int = 30
const maxFileSize int = 2048000

var allowedFileExtensions = []string{
	".jpg",
	".jpeg",
	".gif",
	".png",
}

var fileExtRegex = regexp.MustCompile("^[a-z0-9]{" + string(minHashLen) + "," + string(maxHashLen) + "}\\.?[a-z0-9]{0,5}$")

func loadConfiguration() {
	file, err := os.Open("mangonel-config.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	if configuration.Length > maxHashLen {
		configuration.Length = maxHashLen
	}

	if configuration.Length < minHashLen {
		configuration.Length = minHashLen
	}

	f, err := os.Stat(configuration.Storage)
	if err != nil {
		panic(err)
	}

	if !f.IsDir() {
		panic("configuration.storage path is not a valid directory")
	}
}

func main() {
	loadConfiguration()
	println("Loading routes")

	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/up", uploadHandler).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/static"))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(configuration.Storage)))

	println("Run")
	err := http.ListenAndServe(configuration.Addr, router)
	if err != nil {
		log.Fatal(err)
	}
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
	validFilesize := false

	reader, err := r.MultipartReader()
	if err != nil {
		jsonErr(w, err)
		return
	}

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
			filesize, err := strconv.Atoi(string(b))
			if err != nil {
				jsonErr(w, errors.New("Unable to read file size"))
				return
			}
			if filesize > maxFileSize {
				jsonErr(w, errors.New("File is too big (>"+string(maxFileSize)+")"))
				return
			}
			validFilesize = true
		case "qqfile":
			if !validFilesize {
				jsonErr(w, errors.New("Unable to determine file size"))
				return
			}
			ext := fileExtension(part.FileName(), part.Header.Get("content-type"))
			if ext == "" {
				jsonErr(w, errors.New("Unable to determine file extension"))
				return
			}
			if !stringInSlice(ext, allowedFileExtensions) {
				jsonErr(w, errors.New("Disallowed file extension : "+ext))
				return
			}

			filename := genFilename(ext)
			err := storeFile(path.Join(configuration.Storage, filename), part)
			files = append(files, filename)
			validFilesize = false
			if err != nil {
				jsonErr(w, err)
				return
			}
			println(part.FileName() + " stored as " + filename)
		}
	}

	err = json.NewEncoder(w).Encode(response{
		Success:   true,
		Filenames: files,
	})

	if err != nil {
		println(err)
	}
}

func jsonErr(w http.ResponseWriter, e error) {
	err := json.NewEncoder(w).Encode(response{
		Success: false,
		Error:   e.Error(),
	})

	if err != nil {
		println(err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func genFilename(ext string) string {
	for {
		filename := uniuri.NewLen(configuration.Length) + ext

		_, err := os.Stat(path.Join(configuration.Storage, filename))
		if err != nil {
			return filename
		}
	}
}

func fileExtension(filename, mimetype string) string {
	ext := path.Ext(filename)
	if ext != "" {
		return ext
	}

	exts, err := mime.ExtensionsByType(mimetype)
	if err != nil || len(exts) == 0 {
		return ""
	}

	return exts[0]
}

func storeFile(to string, part *multipart.Part) error {
	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, part); err != nil {
		return err
	}

	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/index.html")
}
