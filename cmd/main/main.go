package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"infinilapse-unified/pkg/cloud"
	"infinilapse-unified/pkg/compiler"
	"infinilapse-unified/pkg/dslrMgmt"
	"infinilapse-unified/pkg/webcamMgmt"
	"os"
	"strconv"
	"time"
)

func main() {
	// Env stuffs
	intTimelapseIntervalMins := getEnvTimelapseInterval()

	s := gocron.NewScheduler(time.UTC)

	if os.Getenv("WEBCAM_CAPTURE") == "false" && os.Getenv("DSLR_CAPTURE") == "false" {
		fmt.Println("Not capturing any cameras.")
	} else {
		_, ierr := s.Every(intTimelapseIntervalMins).Minutes().Do(CaptureAllCameras)
		if ierr != nil {
			println("oh no -- %s", ierr)
		}
	}

	if os.Getenv("COMPILE") == "false" {
		println("NOT COMPILING")
	} else {
		var yesterdayCompileErr error
		if os.Getenv("COMPILE_NOW") == "true" {
			_, yesterdayCompileErr = s.Every(1).Day().Do(compiler.ChunkCompiler)
		} else {
			_, yesterdayCompileErr = s.Every(1).Day().At("00:01").Do(compiler.ChunkCompiler)
		}
		if yesterdayCompileErr != nil {
			println("chunky err --- %s", yesterdayCompileErr)
		}
	}

	//s.StartAsync()
	s.StartBlocking()

}

func CaptureAllCameras() {
	var capturedFiles []string
	if os.Getenv("DSLR_CAPTURE") == "false" {
		println("NOT CAPTURING DSLR")
	} else {
		capturedFiles = append(capturedFiles, dslrMgmt.CaptureAllDslr()...)
		fmt.Printf("dslrMgmt.CaptureAllDslr()...\n%v\n", capturedFiles)
	}
	if os.Getenv("WEBCAM_CAPTURE") == "false" {
		println("NOT CAPTURING WEBCAM")
	} else {
		capturedFiles = append(capturedFiles, webcamMgmt.CaptureWebCams()...)
		fmt.Printf("webcamMgmt.CaptureWebCams()...\n%v\n", capturedFiles)
	}

	err := cloud.IndexGoogleCloudStorageAndGraphQL(capturedFiles)
	if err != nil {
		fmt.Errorf("cloud.IndexGoogleCloudStorageAndGraphQL(filePaths) %s\n", err)
	}

	fmt.Printf("Finished cap loop.\n")
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
