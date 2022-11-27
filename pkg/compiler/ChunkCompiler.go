package compiler

import (
	"errors"
	"fmt"
	"github.com/bitfield/script"
	"infinilapse-unified/pkg/color"
	"infinilapse-unified/pkg/envHelp"
	"infinilapse-unified/pkg/gcpMgmt"
	"infinilapse-unified/pkg/gqlMgmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DEBUG = false
)

func outDirFromInDir(inputDir string) string {
	baseDir := envHelp.BaseDirFromEnv()

	parts := strings.Split(inputDir, "/")
	outSlug := parts[len(parts)-1]
	outDir := baseDir + "/data/out/" + outSlug

	return outDir
}

func ListAvailableDates(camDir string) []string {
	files, err := ioutil.ReadDir(camDir)
	if err != nil {
		log.Fatal(err)
	}

	var miList []string
	daysMap := make(map[string]bool)

	for _, file := range files {
		//fmt.Println(file.Name(), file.IsDir())
		dateStringSlice := strings.Split(file.Name(), ".jpg")
		dateString := dateStringSlice[0]

		dayStringSlice := strings.Split(dateString, "T")

		var dayString string
		if len(dayStringSlice) > 1 {
			dayString = dayStringSlice[0]
			daysMap[dayString] = true
		}
	}

	for key, element := range daysMap {
		fmt.Println("Key:", key, "=>", "Element:", element)
		miList = append(miList, key)
	}

	return miList
}

func CompileMissing(camDir string, available []string) []string {
	outDir := outDirFromInDir(camDir)

	var compiledList []string

	for _, date := range available {
		correspondingOutChunk := outDir + "/" + date + ".mp4"
		if info, err := os.Stat(correspondingOutChunk); err == nil {
			// path/to/whatever exists
			fmt.Printf("%s exists...\n", correspondingOutChunk)
			fmt.Printf("%d bytes large\n", info.Size())

		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			fmt.Printf("%s does not exist... we are going to compile it now.\n", correspondingOutChunk)
			compiledList = append(compiledList, compileDayFromDirAndDate(camDir, date))
		}
	}

	return compiledList
}

func ChunkCompiler() {
	color.PrintCyanBold("begin ChunkCompiler()")

	basedir := envHelp.BaseDirFromEnv()

	inputDir := basedir + "/data/img/dslr"

	cameraDirList := listCameras(inputDir)

	anotherDir := basedir + "/data/img/webcams"

	cameraDirList = append(cameraDirList, listCameras(anotherDir)...)

	println("===  WELCOME TO THE DAILY FFMPEG COMPILER  ===\n")

	fmt.Printf("%v\n", cameraDirList)

	for _, camDir := range cameraDirList {
		availableStillsForCompile := ListAvailableDates(camDir)
		compiledList := CompileMissing(camDir, availableStillsForCompile)

		//outMp4Path := compileDayFromDirAndDate(camDir, dateFromOffset(dateOffsetInt))

		color.PrintCyanBold(fmt.Sprintf("Produced paths: %v", compiledList))

		fmt.Printf("=== INDEX ===\n")
		if os.Getenv("DRY_RUN") == "TRUE" {
			fmt.Printf("DRY_RUN==TRUE .. not indexing")
		} else {
			for _, outMp4Path := range compiledList {
				err := IndexChunk(outMp4Path, lastPartFromPath(camDir), "gcb-site-pub")
				if err != nil {
					fmt.Errorf("IndexChunk(...) --- %s\n", err)
				}
			}
		}

		// manual compile. archival only
		//compileDayFromDirAndDate(fmt.Sprintf("%s/%s", "./sandboxfiles/data/img/dslr", cam), cam)
	}

	// put any outstanding chunks together and update latest.
	_, err := CompileAllPreviousVideo()
	if err != nil {
		_ = fmt.Errorf("error compiling prev: %s\n", err)
		return
	}
}

func IndexChunk(photoFilePath, cameraName, bucket string) error {
	var fileName string
	fileName = lastPartFromPath(photoFilePath)

	nodeName := sanitizeGetEnvNodeName()

	//bucket := "tl-data"
	objPath := "tl-chunk" + "/" + nodeName + "/" + cameraName + "/" + fileName

	fmt.Printf("\n\n==========\nIndexing file at:\n%s to GCP bucket at path: %s\n", photoFilePath, bucket+"/"+objPath)

	objUrl, err := gcpMgmt.StoreFileToBucket(photoFilePath, objPath, bucket)
	if err != nil {
		fmt.Printf("StoreFileToBucket err --- %s\n", err)
	}
	// upload was success index in the graph
	gqlMgmt.IndexToGraph(objUrl, bucket, nodeName+"."+cameraName+".TL-CHUNK")

	return nil
}

func sanitizeGetEnvNodeName() string {
	nodeName := os.Getenv("MY_NODE_NAME")
	if nodeName != "" {
		return nodeName
	} else {
		return "macbook"
	}
}

func lastPartFromPath(photoFilePath string) string {
	parts := strings.Split(photoFilePath, "/")
	return parts[len(parts)-1]
}

func dateFromOffset(offsetDays int) string {
	dateMoment := time.Now().AddDate(0, 0, offsetDays)
	return fmt.Sprint(dateMoment.Format("2006-01-02"))
}

func yesterdayDateString() string {
	return dateFromOffset(-1)
}

func compileDayFromDirAndDate(inputDir string, dayString string) (outMp4PathString string) {

	outDir := outDirFromInDir(inputDir)

	mkdirCmd := "mkdir -p " + outDir

	execStdoutWrapper(mkdirCmd)

	fileName := dayString + ".mp4"
	outMp4PathString = outDir + "/" + fileName

	rmCmd := "rm -f " + outMp4PathString

	execStdoutWrapper(rmCmd)

	compileExecString := fmt.Sprintf(
		"ffmpeg -y -f image2 -r 60 -pattern_type glob -i '%s/%s*.jpg' -vcodec libx264  -pix_fmt yuv420p %s",
		inputDir,
		dayString,
		outMp4PathString,
	)

	execStringWrapper(compileExecString)

	return outMp4PathString
}

func execStringWrapper(cmd string) {
	color.PrintCyanBold(cmd)

	if os.Getenv("DRY_RUN") == "TRUE" {
		fmt.Printf("DRY_RUN set not running cmd:\n%s\n", cmd)
	} else {
		//_, returnCode := script.Exec(compileExecString).Stdout()
		resultStr, returnCode := script.Exec(cmd).String()
		if DEBUG {
			fmt.Printf("%s\n", resultStr)
		}
		if returnCode != nil {
			_ = fmt.Errorf("bad return code `ffmpeg` cmd: %s", returnCode)
		}
	}
}

func execStdoutWrapper(cmd string) {
	color.PrintCyanBold(cmd)
	if os.Getenv("DRY_RUN") == "TRUE" {
		fmt.Printf("DRY_RUN==TRUE .. not running\n")
	} else {
		_, err := script.Exec(cmd).Stdout()
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}

func readMiFiles() {
	files, err := ioutil.ReadDir("./sandboxfiles/data/img/dslr/6D")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		//fmt.Println(file.Name(), file.IsDir())
		dateStringSlice := strings.Split(file.Name(), ".jpg")
		dateString := dateStringSlice[0]
		dateObj, err := time.Parse("2006-01-02T15:04:05+0000", dateString)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d\n", dateObj.Day())
	}
}

func listCameras(inputDir string) (cameraList []string) {

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			cameraList = append(cameraList, inputDir+"/"+file.Name())
		} else {
			fmt.Errorf("/data/img/dslr/%s --- ooops.  Should only contain directories", file.Name())
		}
	}

	return cameraList
}
