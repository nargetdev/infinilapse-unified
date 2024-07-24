package webcamMgmt

import (
	"fmt"
	"github.com/bitfield/script"
	"infinilapse-unified/pkg/envHelp"
	"infinilapse-unified/pkg/gcpMgmt"
	"infinilapse-unified/pkg/gqlMgmt"
	"infinilapse-unified/pkg/util"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func CaptureWebCams() (filePaths []string) {
	devicesList := EnumerateUsbWebCamDevices()
	filePaths = CaptureFromDevicesList(devicesList)

	return filePaths
}

type exposureMap struct {
	exp0 int
	exp2 int
	exp4 int
	exp6 int
	exp8 int
}

func CaptureFromDevicesList(devices []string) (imgFilePathsList []string) {
	// Initialize paths, etc.
	if DebugEnabled {
		fmt.Printf("func CaptureFromDevicesList(devices []string) (synopsis string) {...}\ndevices: %v\n", devices)
	}
	DataDir := os.Getenv("OUTPUT_DIR_WEBCAMS")
	if DataDir == "" {
		DataDir = "/data/" + envHelp.NodeNameFromEnv() + "/img/webcams"
	}

	var exposureValues exposureMap = getExposureMapFromEnv()

	for _, device := range devices {

		deviceLastPart := util.LastPartFromPath(device)
		dataStoreDir := fmt.Sprintf("%s/%s", DataDir, deviceLastPart)
		script.Exec("mkdir -p " + dataStoreDir)

		fileName := fmt.Sprintf("%s.jpg", getDateStringNow())

		fullPathResult := fmt.Sprintf("%s/%s", dataStoreDir, fileName)

		exposureAbsolute := expValueByDevice(exposureValues, deviceLastPart)

		fmt.Printf("Capture on device: %s --- at exposure %d\n", device, exposureAbsolute)

		//		cmdStr := fmt.Sprintf(`
		//v4l2-ctl -d %s \
		//--set-fmt-video=width=1920,height=1080 \
		//--stream-mmap \
		//--stream-count=1 \
		//--stream-to=%s
		//`, device, fullPathResult)
		cmdStr := fmt.Sprintf(`ffmpeg -f v4l2 -video_size 1920x1080 -i %s -frames:v 1 %s`, device, fullPathResult)

		for i := 0; i < 2; i++ {
			outStr, err := script.Exec(cmdStr).String()
			if err != nil {
				_ = fmt.Errorf("v4l2 error --- %s\n", err)
			}
			if DebugEnabled {
				fmt.Println(outStr)
			}
		}

		imgFilePathsList = append(imgFilePathsList, fullPathResult)
	}

	return imgFilePathsList
}

func expValueByDevice(values exposureMap, deviceStrShort string) int {
	switch deviceStrShort {
	case "video0":
		return values.exp0
	case "video2":
		return values.exp2
	case "video4":
		return values.exp4
	case "video6":
		return values.exp6
	case "video8":
		return values.exp8
	}
	_ = fmt.Errorf("Oops. We did not select an exposure. %s\n", deviceStrShort)
	return 1234 // shouldn't ever happen.
}

func getExposureMapFromEnv() exposureMap {
	v0, _ := strconv.Atoi(os.Getenv("EXPOSURE_0"))
	v2, _ := strconv.Atoi(os.Getenv("EXPOSURE_2"))
	v4, _ := strconv.Atoi(os.Getenv("EXPOSURE_4"))
	v6, _ := strconv.Atoi(os.Getenv("EXPOSURE_6"))
	v8, _ := strconv.Atoi(os.Getenv("EXPOSURE_8"))

	var myExposures = exposureMap{
		exp0: v0,
		exp2: v2,
		exp4: v4,
		exp6: v6,
		exp8: v8,
	}

	return myExposures
}

func ExecCamCap02468_fswebcam() {
	DataDir := os.Getenv("OUTPUT_DIR_WEBCAMS")
	if DataDir == "" {
		DataDir = "/data/img/webcams"
	}

	staticPathTestUpload := os.Getenv("STATIC_PATH_TEST_UPLOAD")
	nodeName := os.Getenv("MY_NODE_NAME")

	fmt.Printf("======================\n")
	fmt.Printf(
		"node: %s\tpod: %s\tpodns: %s\tpodip: %s",
		nodeName,
		os.Getenv("MY_POD_NAME"),
		os.Getenv("MY_POD_NAMESPACE"),
		os.Getenv("MY_POD_IP"),
	)
	fmt.Printf("\n======================\n")

	shSelectStr := "sh"

	dateString := getDateStringNow()

	cameras := [5]int{0, 2, 4, 6, 8}

	for i := 0; i < len(cameras); i++ {
		vidId := fmt.Sprintf("video%d", cameras[i])

		dataStoreDir := fmt.Sprintf("%s/%s", DataDir, vidId)

		cmdStr := fmt.Sprintf("mkdir -p %s", dataStoreDir)

		// TODO: refactor these exec.Command to it's own function
		_, err := exec.Command(shSelectStr, "-c", cmdStr).Output()
		if err != nil {
			log.Println(err)
		}

		fileName := fmt.Sprintf("%s.jpg", dateString)

		fullPath := fmt.Sprintf("%s/%s", dataStoreDir, fileName)

		fsWebCamCmd := fmt.Sprintf("fswebcam -r 1920x1080 --device /dev/%s %s", vidId, fullPath)
		//fsWebCamCmd := fmt.Sprintf("fswebcam --device /dev/%s %s", vidId, fullPath)

		println(fsWebCamCmd)

		fsWebCamOutput, err := exec.Command(shSelectStr, "-c", fsWebCamCmd).CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		log.Printf("%s", fsWebCamOutput)

		//bucket := "tl-data"
		bucket := os.Getenv("CLOUD_STORAGE_BUCKET")
		if bucket == "" {
			bucket = "gcb-site-pub"
			println("\nno CLOUD_STORAGE_BUCKET env --- setting static: " + bucket)
		}

		objPath := nodeName + "/" + vidId + "/" + fileName

		// allow time for capture to complete
		time.Sleep(time.Second * 2)
		// check if file size is zero (command not returned yet)
		//for nUtils.FileIsEmpty(fullPath) {
		//	fmt.Printf("%s is still empty", fullPath)
		//}

		if staticPathTestUpload != "" {
			fmt.Printf("Have static path, passing it:\n%s\n", staticPathTestUpload)
			fullPath = staticPathTestUpload
		}
		//UploadFile
		errr := gcpMgmt.UploadFile(fullPath, objPath, bucket)
		if errr != nil {
			fmt.Printf("%s", errr)
			continue
		}

		objUrl := "https://storage.googleapis.com/" + bucket + "/" + objPath

		// upload was success index in the graph
		gqlMgmt.IndexToGraph(objUrl, bucket, nodeName+"."+vidId)
	}
	println("===========================")
	println("===========================")
	println("===========================")
}

func getDateStringNow() string {
	tm := time.Now()
	dateString := tm.Format(time.RFC3339)
	return dateString
}
