package webcamMgmt

import (
	"bufio"
	"fmt"
	"github.com/bitfield/script"
	"regexp"
	"strings"
)

const (
	DebugEnabled = true
)

// EnumerateUsbWebCamDevices
// Depends:
// v4l2-ctl
func EnumerateUsbWebCamDevices() []string {
	cmdStr := "v4l2-ctl --list-devices"

	outStr, _ := script.Exec(cmdStr).String()

	fmt.Println(outStr)

	return DevicesStringListFromListDevices(outStr)
}

func DevicesStringListFromListDevices(rawString string) (webcamDevices []string) {

	var wantedDeviceOnNextLine bool = false
	pattern := regexp.MustCompile("\\(usb-.*\\)")
	scanner := bufio.NewScanner(strings.NewReader(rawString))
	for scanner.Scan() {
		thisLineString := scanner.Text()
		if wantedDeviceOnNextLine {
			webcamDevices = append(webcamDevices, strings.TrimSpace(thisLineString))
			wantedDeviceOnNextLine = false
			continue
		}

		firstMatchIndex := pattern.FindStringIndex(thisLineString)
		if firstMatchIndex != nil {
			fmt.Println("First matched index", firstMatchIndex[0], "-", firstMatchIndex[1])
			fmt.Println(getSubstring(thisLineString, firstMatchIndex))
			fmt.Println(scanner.Text())
			wantedDeviceOnNextLine = true
		}
		fmt.Println()
	}

	return webcamDevices
}

func getSubstring(s string, indices []int) string {
	return string(s[indices[0]:indices[1]])
}
