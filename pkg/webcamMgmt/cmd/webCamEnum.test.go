package main

import (
	"infinilapse-unified/pkg/webcamMgmt"
)

func main() {
	devicesList := webcamMgmt.EnumerateUsbWebCamDevices()
	webcamMgmt.CaptureFromDevicesList(devicesList)

	// sleep now kill container manually
	select {}
}
