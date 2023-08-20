package core

import (
	"fmt"
	"runtime"
)

func TcpError(e string) {
	fmt.Println("[tcpError] ", e)
}

func CoreInfo(i string) {
	fmt.Println("[coreInfo]", i)
}

func CoreError(err error) {
	// notice that we're using 1, so it will actually log where
	// the error happened, 0 = this function, we don't want that.
	_, filename, line, _ := runtime.Caller(1)
	fmt.Printf("[core error] %s:%d %v", filename, line, err)
}
