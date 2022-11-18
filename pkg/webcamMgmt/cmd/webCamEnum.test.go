package main

import (
	"fmt"
	"infinilapse-unified/pkg/webcamMgmt"
)

func main() {
	fmt.Println("ohai")
	webcamMgmt.EnumerateUsbWebCams()
}
