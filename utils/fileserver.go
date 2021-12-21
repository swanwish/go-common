package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ServeDownloadFile(filePath string, rw http.ResponseWriter, r *http.Request) {
	serveFile(rw, r, filePath, true)
}

func ServeFile(rw http.ResponseWriter, r *http.Request, filePath string) {
	serveFile(rw, r, filePath, false)
}

func serveFile(rw http.ResponseWriter, r *http.Request, filePath string, download bool) {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(rw, "404 Not Found : Error while opening the file.", 404)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	stat, err := file.Stat()
	if err != nil {
		http.Error(rw, "500 Internal Error : stat() failure.", 500)
		return
	}

	if download {
		_, fileName := filepath.Split(filePath)
		if fileName != "" {
			rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		}
	}

	http.ServeContent(rw, r, stat.Name(), stat.ModTime(), file)
}

func ServeContent(rw http.ResponseWriter, r *http.Request, name string, content []byte) {
	serveContent(rw, r, name, content, time.Now(), false)
}

func ServeDownloadContent(rw http.ResponseWriter, r *http.Request, name string, content []byte) {
	serveContent(rw, r, name, content, time.Now(), true)
}

func serveContent(rw http.ResponseWriter, r *http.Request, name string, content []byte, modtime time.Time, download bool) {
	reader := bytes.NewReader(content)

	if download && name != "" {
		rw.Header().Set("Content-Disposition", name)
	}

	http.ServeContent(rw, r, name, modtime, reader)
}
