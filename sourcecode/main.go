package main

import (
	"gioms/server"
)

func main() {
	var server = server.NewServer("localhost:123")
	server.Start()
}
