package main

import (
	"peppa.pig.com/model_test/server"
)

func main() {
	newServer := server.NewServer("127.0.0.1", 8080)

	newServer.Start()
}
