package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/dchest/uniuri"
)

func genFilename(ext string) string {
	for {
		filename := uniuri.NewLen(configuration.HashLength) + ext
		log.Debugf("Generated filename : %s", filename)

		if !fileExists(path.Join(configuration.StoragePath, filename)) {
			return filename
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func fileExtension(filename, mimeType string) string {
	ext := path.Ext(filename)
	if ext != "" {
		return strings.ToLower(ext)
	}

	extArray, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(extArray) == 0 {
		return ""
	}

	return extArray[0]
}

func handleFilePart(part *multipart.Part) (string, error) {
	log.Debugf("Extract extension from %s", part.FileName())
	ext := fileExtension(part.FileName(), part.Header.Get("content-type"))
	if ext == "" || !stringInSlice(ext, allowedFileExtensions) {
		return "", errors.New("unknown file type")
	}

	b, err := ioutil.ReadAll(part)
	if err != nil {
		return "", err
	}

	r := bytes.NewReader(b)
	h := sha256.New()
	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))

	log.Debugf("File hashed to %s", hash)
	hashPath := path.Join(configuration.StoragePath, hash)
	if !fileExists(hashPath) {
		_, err = r.Seek(0, io.SeekStart)
		if err != nil {
			return "", err
		}
		err = storeFile(hashPath, r)
		if err != nil {
			return "", err
		}
	}

	filename := genFilename(ext)
	err = os.Link(hashPath, path.Join(configuration.StoragePath, filename))
	if err != nil {
		return "", err
	}

	log.Debugf("Stored as %s", filename)
	return filename, nil
}

func storeFile(to string, src io.Reader) error {
	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
