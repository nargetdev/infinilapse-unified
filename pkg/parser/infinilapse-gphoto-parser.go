package parser

import (
	"fmt"
	"regexp"
	"strings"
)

func NamesAndPortsFromMultiLineAutoDetect(autoDetectMultilineString string) [][]string {
	gphotoAutoDetectCameraList := strings.Split(autoDetectMultilineString, "\n")
	fmt.Printf("%#v\n\n", gphotoAutoDetectCameraList)

	var namesAndPortsArray [][]string
	for idx, autoDetectLineString := range gphotoAutoDetectCameraList {
		if autoDetectLineString != "" {
			//nameAndPortArrayElt := strings.Split(autoDetectLineString, "    ") // legacy.. and broken hah
			model, port := ParseOneLineGphotoAutoDetect(autoDetectLineString)
			namesAndPortsArrayElt := []string{model, port}
			fmt.Printf("%d --- %#v\n\n", idx, namesAndPortsArrayElt)
			namesAndPortsArray = append(namesAndPortsArray, namesAndPortsArrayElt[0:2])
		}
	}

	return namesAndPortsArray
}
func ParseOneLineGphotoAutoDetect(sampleInput string) (model, port string) {
	r, _ := regexp.Compile("usb:\\d\\d\\d,\\d\\d\\d")

	//println(r.FindString(sampleInput))

	portIndices := r.FindStringIndex(sampleInput)
	//fmt.Println(portIndices)
	//println(portIndices)

	port = sampleInput[portIndices[0]:portIndices[1]]
	//println(port)

	model = sampleInput[0:portIndices[0]]
	//println(model)
	model = strings.TrimSpace(model)

	return model, port
}
