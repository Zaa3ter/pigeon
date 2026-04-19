package main

import (
	"fmt"
	"os"

	"github.com/MJ-NMR/pigeon/client"
	"github.com/MJ-NMR/pigeon/server"
)

func main() {
	if len(os.Args) < 3 {
		help()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "-c":
		client.Connect(os.Args[2])
	case "-l":
		server.Listen(os.Args[2])
	}
}

func help() {
	fmt.Println("server: pigeon -l <port>")
	fmt.Println("client: pigeon -c <address>:<port>")
}
