package server

import (
	"net"
)

type playerSession struct {
	conn    net.Conn
	Version int
}
