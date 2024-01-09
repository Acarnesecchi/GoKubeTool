package fileserver

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	config Config
}

func NewServer(c Config) Server {
	return Server{config: c}
}

func StartServer() {
	c, err := NewConfig().WithStorage("local", "/opt/devtools/storage2")
	if err != nil {
		os.Exit(1)
	}
	s := NewServer(c)

	http.HandleFunc("/", serveForm)
	http.HandleFunc("/upload", handleFileUpload)
	fmt.Printf("Server starting at http://localhost%s", s.config.listenAddr)
	http.ListenAndServe(s.config.listenAddr, nil)
}
