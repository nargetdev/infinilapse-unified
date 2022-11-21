package dslrMgmt

import (
	"fmt"
	"github.com/bitfield/script"
	"infinilapse-unified/pkg/parser"
	"strings"
)

func CaptureAllDslr() (files []string) {
	//GcbBucket := "gcb-site-pub"
	// Get gphoto2 autodetect.  First line of result will be the "----" line.
	autoDetectMultilineString, _ := script.Exec("gphoto2 --auto-detect").Exec("tail -n +3").String()

	namesAndPortsArray := parser.NamesAndPortsFromMultiLineAutoDetect(autoDetectMultilineString)

	fmt.Println("OUR CAMERAS:")
	fmt.Printf("%#v\n\n", namesAndPortsArray)

	for _, cam := range namesAndPortsArray {
		//fmt.Printf("\n=====loop(namesAndPortsArray)===\nAbout to execute from array... %#v\n\n", namesAndPortsArray)

		cameraName := strings.Split(cam[0], " ")[2]

		ImgDir := ""
		ImgDir = ImgDir + "/data/img/dslr/" + cameraName

		datetimeStringWithNewline, err := script.Exec("date +%Y-%m-%dT%H:%M:%S%z").String()
		datetimeString := strings.TrimSuffix(datetimeStringWithNewline, "\n")
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		fileName := fmt.Sprintf("%s.jpg", datetimeString)
		FullPath := fmt.Sprintf("%s/%s", ImgDir, fileName)
		//captureExecString := fmt.Sprintf("gphoto2 --debug --debug-loglevel debug --camera \"%s\" --port \"%s\" --capture-image-and-download --filename \"%%y-%%m-%%d__%%H:%%M:%%S.jpg\"", cam[0], cam[1])
		captureExecString := fmt.Sprintf("timeout 16 gphoto2 --camera \"%s\" --port \"%s\" --capture-image-and-download --filename \"%s\"", cam[0], cam[1], FullPath)
		fmt.Printf("EXECUTING> %s\n", captureExecString)
		bytesWrote, returnCode := script.Exec(captureExecString).Stdout()
		if returnCode != nil {
			fmt.Printf("STATUS:: %s\n continuing to next camera...\n", returnCode)
			continue
		}
		fmt.Printf("Wrote %d bytes and received exit status :: %s\n\n", bytesWrote, returnCode)

		files = append(files, FullPath)

		//err = cloud.IndexPhoto(FullPath, cameraName, fileName, GcbBucket)
		//if err != nil {
		//	fmt.Printf("IndexPhoto(...)  --- err != nil {\n%s\n}", err)
		//}
	}
	return files
}
