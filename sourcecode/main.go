package main

import (
  "gioms/server"
  "gioms/utils"
  "fmt"
)

func main() {
  var server := NewServer("localhost:123")
  server.Start()
}