package server

import(
	"net"
	"core"
)

type struct MinecraftServer {
	listenAddr string
	ln net.Listener
	quitch chan struct{}
}

func NewServer(listenAddr string) *MinecraftServer {
	return &MinecraftServer{
		listenAddr: listenAddr,
		quitch: make(chan struct{})
	}
}

func (s *MinecraftServer) Start() error {
	ln, err := net.Listener("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	<-quitch

	return nil
}

func(s *MinecraftServer) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
	}
}