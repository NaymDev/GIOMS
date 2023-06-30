package core

import(
	"fmt"
)

func TcpError(e string) {
	fmt.Println("[tcpError] ", e)
}

func CoreInfo(i string) {
	fmt.Println("[coreInfo]", i)
}

func CoreError(e error) {
	fmt.Println("[coreError]", e)
}