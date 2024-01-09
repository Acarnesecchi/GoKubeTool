package fileserver

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func serveForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("upload-form").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>File Upload</title>
		<script src="https://unpkg.com/htmx.org"></script>
	</head>
	<body>
		<h2>Upload File</h2>
		<form id='form' hx-encoding='multipart/form-data' hx-post='/upload'>
			<input type='file' name='file'>
			<button>
				Upload
			</button>
			<progress id='progress' value='0' max='100'></progress>
		</form>
		<script>
			htmx.on('#form', 'htmx:xhr:progress', function(evt) {
			htmx.find('#progress').setAttribute('value', evt.detail.loaded/evt.detail.total * 100)
			});
		</script>
			<div id="upload-status"></div>
	</body>
	</html>
    `))
	tmpl.Execute(w, nil)
}

func createCompressedFile(d, n string, f multipart.File) error {
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error")
	}
	tarFile, err := os.Create(d)
	if err != nil {
		return fmt.Errorf("error")
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	header := &tar.Header{
		Name: n,
		Size: int64(len(fileBytes)),
		Mode: 0600,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("error")
	}
	if _, err := tarWriter.Write(fileBytes); err != nil {
		return fmt.Errorf("error")
	}
	return nil
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		r.ParseMultipartForm(10 << 20) // 10 MB max memory
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		base, _ := os.UserHomeDir()
		destPath := filepath.Join(base, ".gkt/storage/mapsweb", handler.Filename+".tar.gz")
		err = createCompressedFile(destPath, handler.Filename, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "<div>File '%s' compressed and uploaded successfully!</div>", handler.Filename)
	}
}
