package main

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/go-co-op/gocron"
	"infinilapse-unified/pkg/compiler"
	"infinilapse-unified/pkg/gcpMgmt"
	"infinilapse-unified/pkg/gqlMgmt"
	"infinilapse-unified/pkg/parser"
	"infinilapse-unified/pkg/webcamMgmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Env stuffs
	intTimelapseIntervalMins := getEnvTimelapseInterval()

	s := gocron.NewScheduler(time.UTC)

	_, ierr := s.Every(intTimelapseIntervalMins).Minutes().Do(CaptureAllCameras)
	if ierr != nil {
		println("oh no -- %s", ierr)
	}

	_, webcamErr := s.Every(intTimelapseIntervalMins).Minutes().Do(CaptureWebCams)
	if webcamErr != nil {
		println("webcam issues --- %s", webcamErr)
	}

	_, yesterdayCompileErr := s.Every(intTimelapseIntervalMins).Day().Do(compiler.ChunkCompiler)
	if yesterdayCompileErr != nil {
		println("webcam issues --- %s", yesterdayCompileErr)
	}

	_, compileErr := s.Every(intTimelapseIntervalMins).Day().Do(compiler.CompileAllPreviousVideo)
	if compileErr != nil {
		println(".webcam issues --- %s", compileErr)
	}

	//s.StartAsync()
	s.StartBlocking()
}

func CaptureWebCams() {
	webcamMgmt.ExecCamCap02468()
}

func CaptureAllCameras() {
	GcbBucket := "gcb-site-pub"
	// Get gphoto2 autodetect.  First line of result will be the "----" line.
	autoDetectMultilineString, _ := script.Exec("gphoto2 --auto-detect").Exec("tail -n +3").String()

	namesAndPortsArray := parser.NamesAndPortsFromMultiLineAutoDetect(autoDetectMultilineString)

	fmt.Println("OUR CAMERAS:")
	fmt.Printf("%#v\n\n", namesAndPortsArray)

	for _, cam := range namesAndPortsArray {
		fmt.Printf("\n=====loop(namesAndPortsArray)===\nAbout to execute from array... %#v\n\n", namesAndPortsArray)

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

		err = IndexPhoto(FullPath, cameraName, fileName, GcbBucket)
		if err != nil {
			fmt.Printf("IndexPhoto(...)  --- err != nil {\n%s\n}", err)
		}
	}

	println("finished cap loop")
}

func IndexPhoto(photoFilePath, cameraName, fileName, bucket string) error {
	fmt.Printf("Indexing photo to GCP and GQL at path: %s\n", photoFilePath)

	nodeName := os.Getenv("MY_NODE_NAME")

	//bucket := "tl-data"
	objPath := nodeName + "/" + cameraName + "/" + fileName
	objUrl, err := gcpMgmt.StoreFileToBucket(photoFilePath, objPath, bucket)
	if err != nil {
		fmt.Printf("StoreFileToBucket err --- %s\n", err)
	}
	// upload was success index in the graph
	gqlMgmt.IndexToGraph(objUrl, bucket, nodeName+"."+cameraName)

	return nil
}

func getEnvTimelapseInterval() interface{} {
	TimelapseIntervalMins := os.Getenv("TIMELAPSE_INTERVAL_MINS")
	fmt.Printf("TIMELAPSE_INTERVAL_MINS: %s\n", TimelapseIntervalMins)
	var intTimelapseIntervalMins int
	var err error
	if TimelapseIntervalMins != "" {
		var atoiErr error
		intTimelapseIntervalMins, atoiErr = strconv.Atoi(TimelapseIntervalMins)
		if atoiErr != nil {
			fmt.Printf("Error converting TIMELAPSE_INTERVAL_MINS to int --- %s", atoiErr)
		}
	} else {
		intTimelapseIntervalMins = 15
	}
	if err != nil {
		fmt.Printf("%s", err)
	}
	return intTimelapseIntervalMins
}
