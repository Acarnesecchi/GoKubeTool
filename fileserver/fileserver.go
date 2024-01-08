package fileserver

import (
	"fmt"
	"os"
)

type Server struct {
	config Config
}

func NewServer(c Config) Server {
	return Server{config: c}
}

func startServer() {
	c, err := NewConfig().WithStorage("local", "/opt/devtools/storage2")
	if err != nil {
		os.Exit(1)
	}
	s := NewServer(c)
	fmt.Printf("server %v", s)
}
