package compiler

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/fatih/color"
	"infinilapse-unified/pkg/gcpMgmt"
	"infinilapse-unified/pkg/gqlMgmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func PrintYellow(say string) {
	yellow := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, say)
	fmt.Println(yellow)
}

func PrintMagenta(say string) {
	cyan := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 35, say)
	fmt.Println(cyan)
}

func PrintCyanBold(say string) {
	cyan := fmt.Sprintf("\x1b[1m\x1b[%dm%s\x1b[0m\x1b[22m", 36, say)
	fmt.Println(cyan)
}
func PrintCyan(say string) {
	cyan := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 36, say)
	fmt.Println(cyan)
}

func ChunkCompiler() {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, "begin ChunkCompiler()")

	PrintYellow("chunkyYellow")
	PrintCyan("ohai cyan")
	PrintMagenta("begin ChunkCompiler()")

	fmt.Println(colored)
	//readMiFiles()
	color.Red("We have red")
	color.Magenta("And many others ..")

	inputDir := "/data/img/dslr"

	cameraDirList := listCameras(inputDir)

	anotherDir := "/data/img/webcams"

	cameraDirList = append(cameraDirList, listCameras(anotherDir)...)

	println("WELCOME TO THE DAILY FFMPEG COMPILER\n")

	fmt.Printf("%v\n", cameraDirList)

	println("===========")
	println("===========")
	println("===========")
	// Use handy standard colors
	color.Set(color.FgYellow)
	// Print with default helper functions
	color.Cyan("Prints text in cyan.")

	// A newline will be appended automatically
	color.Blue("Prints %s in blue.", "text")

	// These are using the default foreground colors
	color.Red("We have red")
	color.Magenta("And many others ..")

	color.Unset()

	DateOffset := os.Getenv("DATE_OFFSET_TO_COMPILE")
	var dateOffsetInt int
	var strconvErr error
	if DateOffset == "" {
		dateOffsetInt = -1
	} else {
		dateOffsetInt, strconvErr = strconv.Atoi(DateOffset)
		if strconvErr != nil {
			fmt.Printf("DATE_OFFSET_TO_COMPILE not converted to int --- %s\n", strconvErr)
		}
	}
	fmt.Printf("\nCONFIG:\n%s\t%d\n", DateOffset, dateOffsetInt)
	for _, camDir := range cameraDirList {
		outMp4Path := compileDayFromDirAndDate(camDir, dateOffsetInt)

		color.Cyan("Produced path: %s", outMp4Path)

		err := IndexChunk(outMp4Path, lastPartFromPath(camDir), "gcb-site-pub")
		if err != nil {
			fmt.Errorf("IndexChunk(...) --- %s\n", err)
		}
		//compileDayFromDirAndDate(fmt.Sprintf("%s/%s", "./sandboxfiles/data/img/dslr", cam), cam)
	}

	// We've finished our tasks.  This container will now close until the next time the control pane does CronJob.

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

	nodeName := os.Getenv("MY_NODE_NAME")

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

func compileDayFromDirAndDate(inputDir string, dateOffset int) (outMp4PathString string) {
	environment := os.Getenv("ENVIRONMENT")
	var baseDir string
	if environment == "dev" {
		baseDir = "."
	} else {
		baseDir = ""
	}

	parts := strings.Split(inputDir, "/")
	outSlug := parts[len(parts)-1]
	outDir := baseDir + "/data/out/" + outSlug
	mkdirCmd := "mkdir -p " + outDir
	_, err := script.Exec(mkdirCmd).Stdout()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	fileName := dateFromOffset(dateOffset) + ".mp4"
	outMp4PathString = outDir + "/" + fileName

	rmCmd := "rm -f " + outMp4PathString
	_, rmErr := script.Exec(rmCmd).Stdout()
	if rmErr != nil {
		fmt.Printf("rm -f ... :::ERR::: %s\n", rmErr)
	}

	compileExecString := fmt.Sprintf(
		"ffmpeg -y -f image2 -r 60 -pattern_type glob -i '%s/%s*.jpg' -vcodec libx264  -pix_fmt yuv420p %s",
		inputDir,
		dateFromOffset(dateOffset),
		outMp4PathString,
	)

	//bytesOfStdout, returnCode := script.Exec(compileExecString).Stdout()

	fmt.Println("our command:")
	fmt.Println(compileExecString)

	if os.Getenv("DRY_RUN") == "yes" {
		fmt.Printf("DRY_RUN set not running cmd:\n%s\n", compileExecString)
	} else {
		_, returnCode := script.Exec(compileExecString).Stdout()
		if returnCode != nil {
			_ = fmt.Errorf("bad return code `ffmpeg` cmd: %s", returnCode)
		}
	}

	return outMp4PathString
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
