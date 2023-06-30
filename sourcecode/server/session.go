package server

import(
	"net"
)

type struct playerSession {
	conn: net.Conn
}