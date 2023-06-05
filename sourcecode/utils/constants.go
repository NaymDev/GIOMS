package utils

type state int8
const (
	handshake = iota
	status 
	login
	play
)