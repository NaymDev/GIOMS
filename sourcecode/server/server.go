package server

import (
	"fmt"
	"gioms/core"
	"gioms/utils"
	"net"
	"strings"
)

type MinecraftServer struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
}

func NewServer(listenAddr string) *MinecraftServer {
	return &MinecraftServer{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *MinecraftServer) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	s.acceptLoop()

	<-s.quitch

	return nil
}

func (s *MinecraftServer) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			core.TcpError("accept error: " + fmt.Sprint(err))
			continue
		}
		sess := playerSession{
			conn: conn,
		}
		go s.doHandshakeWith(sess)
	}
}

func (s *MinecraftServer) doHandshakeWith(session playerSession) {
	defer session.conn.Close()
	buf := make([]byte, 1024)

	n, err := session.conn.Read(buf)
	if handleIfError(err) {
		return
	}

	fmt.Println(buf[:n])
	var rPacket = core.NewServerboundMinecraftPacket()
	rPacket.SetPacked(buf[:n])

	err = rPacket.Unpack("handshake")
	if handleIfError(err) {
		return
	}

	fmt.Println(rPacket.Fields)

	if rPacket.Fields["next_state"] == utils.STATUS {
		//s.sendStatus(session)
		for {
			n, err := session.conn.Read(buf)
			if handleIfError(err) {
				return
			}
			rPacket.SetPacked(buf[:n])

			err = rPacket.Unpack("status")
			if handleIfError(err) {
				return
			}

			if rPacket.RawInfo.PacketID == 0 { //Status
				s.sendStatus(session)
			} else if rPacket.RawInfo.PacketID == 1 { //Ping
				session.conn.Write(buf[:n])
			}
		}
	}
}

func (s *MinecraftServer) sendStatus(session playerSession) {
	var wPacket, err = core.ClientboundPacketWithFields("status", 0)
	handleIfError(err)

	wPacket.RawInfo.State = "status"
	wPacket.Fields["json_response"] = strings.ReplaceAll(strings.ReplaceAll(`{
			"version": {
				"name": "1.19.4",
				"protocol": 762
			},
			"players": {
				"max": 100,
				"online": 5,
				"sample": [
					{
						"name": "thinkofdeath",
						"id": "4566e69f-c907-48ee-8d71-d7ba5aa00d20"
					}
				]
			},
			"description": {
				"text": "Hello world"
			}
		}`, "\n", ""), "	", "")

	fmt.Println("data", string(wPacket.Pack()))
	session.conn.Write(wPacket.Pack())
}

func handleIfError(err error) bool {
	if err != nil {
		core.CoreError(err)
		return true
	}
	return false
}
