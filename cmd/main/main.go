package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
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

	_, ierr := s.Every(intTimelapseIntervalMins).Minutes().Do(CaptureAllCameras)
	if ierr != nil {
		println("oh no -- %s", ierr)
	}

	if os.Getenv("COMPILE") == "false" {
		println("NOT COMPILING")
	} else {
		_, yesterdayCompileErr := s.Every(1).Day().Do(compiler.ChunkCompiler)
		if yesterdayCompileErr != nil {
			println("chunky err --- %s", yesterdayCompileErr)
		}

		// TODO: roll into prev
		_, compileErr := s.Every(1).Day().Do(compiler.CompileAllPreviousVideo)
		if compileErr != nil {
			println("compile err --- %s", compileErr)
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
	}
	if os.Getenv("WEBCAM_CAPTURE") == "false" {
		println("NOT CAPTURING WEBCAM")
	} else {
		capturedFiles = append(capturedFiles, webcamMgmt.CaptureWebCams()...)
	}

	fmt.Printf("Finished cap loop.  Got files:\n%v\n", capturedFiles)
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
