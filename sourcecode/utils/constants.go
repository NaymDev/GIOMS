package utils

type state int8

const (
	HANDSHALE = 0
	STATUS    = 1
	LOGIN     = 2
	PLAY
)

const DEFAULT_PACKET_VERSION = 762
