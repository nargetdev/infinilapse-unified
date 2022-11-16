package main

import (
	"fmt"
	"github.com/bitfield/script"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	readMiFilesAndCopyEmToTheRightSpot()
}

func parseEpochStringtoDateString(epochString string) (dateString string) {
	i, err := strconv.ParseInt(epochString, 10, 64)
	if err != nil {
		fmt.Printf("%s\n", err)
		return ""
	}
	tm := time.Unix(i, 0)
	//dateString = tm.String()
	dateString = tm.Format(time.RFC3339)
	return dateString
}

func readMiFilesAndCopyEmToTheRightSpot() {
	BaseDir := os.Getenv("DATA_BASE_DIR")
	if BaseDir == "" {
		println("cool it's empty")
	} else {
		println("damn it's: " + BaseDir)
	}
	//BaseDir := "."
	//videoDirs := []string{"video0", "video2", "video4", "video6", "video8"}
	videoDirs := []string{"dslr/6D", "dslr/M50"}

	for _, viddir := range videoDirs {
		thisDir := BaseDir + "/data/img/" + viddir
		files, err := ioutil.ReadDir(thisDir)
		if err != nil {
			fmt.Printf("We had an error to read directory %s --- %s\n\n\n::CONTINUE::\n\n", thisDir, err)
			continue
		}

		//outDir := BaseDir + "/data/img/webcams/" + viddir
		outDir := thisDir

		mkdirCmd := "mkdir -p " + outDir
		_, mkdirErr := script.Exec(mkdirCmd).Stdout()
		if mkdirErr != nil {
			fmt.Printf("::ERR:: mkdirCmd\n%s\n", mkdirErr)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			//fmt.Println(file.Name(), file.IsDir())
			dateStringSlice := strings.Split(file.Name(), ".jpg")

			epochTimeString := strings.Trim(file.Name(), ".jpgpic__")
			oldEpochTimeString := dateStringSlice[0]

			if epochTimeString != oldEpochTimeString {
				fmt.Printf("yeaaaa.... %s\n", file.Name())
			}

			dateString := parseEpochStringtoDateString(epochTimeString)
			if dateString == "" {
				continue
			}
			//fmt.Printf("Input epochTimeString: %s ==> Output dateString: %s\n", epochTimeString, dateString)

			// cp file outdir
			cpCmdString := "cp " + thisDir + "/" + file.Name() + " " + outDir + "/" + dateString + ".jpg"

			println(cpCmdString)

			// copy the file to new home.
			_, cpCmdErr := script.Exec(cpCmdString).Stdout()
			if cpCmdErr != nil {
				fmt.Printf("::ERR:: cpCmd\n%s\n", mkdirErr)
			}
		}
	}

}
