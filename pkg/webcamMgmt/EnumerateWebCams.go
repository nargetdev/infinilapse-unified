package webcamMgmt

import (
	"fmt"
	"github.com/bitfield/script"
)

func EnumerateUsbWebCams() {
	cmdStr := "v4l2-ctl --list-devices"

	outStr, _ := script.Exec(cmdStr).String()

	fmt.Println(outStr)

	//autoDetectMultilineString, _ := script.Exec("gphoto2 --auto-detect").Exec("tail -n +3").String()

}
