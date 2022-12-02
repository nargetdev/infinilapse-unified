package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"gitlab.com/jamiemint4096/actionGroupREST"
	"infinilapse-unified/pkg/cloud"
	"infinilapse-unified/pkg/compiler"
	"infinilapse-unified/pkg/dslrMgmt"
	"infinilapse-unified/pkg/oscmanager"
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
	fmt.Printf("Begin cap loop.  Setting the stage\n")
	SetTheStage()

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

	fmt.Printf("Finished cap loop.  Unsetting the stage\n")
	UnsetTheStage()
}

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}

func STOPGAP_getFadeByTime() float64 {
	newLayout := "15:04"
	loc, _ := time.LoadLocation("MST")
	check := time.Now().In(loc)
	start, _ := time.ParseInLocation(newLayout, "08:00", loc)
	end, _ := time.ParseInLocation(newLayout, "23:59", loc)
	start = start.AddDate(check.Year(), int(check.Month())-1, check.Day()-1)
	end = end.AddDate(check.Year(), int(check.Month())-1, check.Day()-1)
	isDaytime := inTimeSpan(start, end, check)
	fmt.Println(start.String()+" --- "+check.String()+" --- ", end.String(), isDaytime)

	if isDaytime {
		return 1.0
	} else {
		return 0.0
	}
}

func UnsetTheStage() {
	actionGroupREST.PostOnOff(actionGroupREST.UUID_fan_relay_0, true, actionGroupREST.PlantiAuthPair)
	brightnessVal := STOPGAP_getFadeByTime()
	oscmanager.FadeMaster(brightnessVal)
}

func SetTheStage() {
	actionGroupREST.PostOnOff(actionGroupREST.UUID_fan_relay_0, false, actionGroupREST.PlantiAuthPair)
	oscmanager.FadeMaster(1.0)
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
