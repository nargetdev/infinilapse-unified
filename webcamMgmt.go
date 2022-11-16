package infinilapse_unified

import (
	"fmt"
	"infinilapse-unified/pkg/gcpMgmt"
	"infinilapse-unified/pkg/gqlMgmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func ExecCamCap02468() {
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

	tm := time.Now()
	dateString := tm.Format(time.RFC3339)

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
		urlResult, errr := gcpMgmt.StoreFileToBucket(fullPath, objPath, bucket)
		if errr != nil {
			fmt.Printf("%s", errr)
			continue
		}

		fmt.Printf(urlResult)

		objUrl := "https://storage.googleapis.com/" + bucket + "/" + objPath

		// upload was success index in the graph
		gqlMgmt.IndexToGraph(objUrl, bucket, nodeName+"."+vidId)
	}
	println("===========================")
	println("===========================")
	println("===========================")
}

//func onMessageReceived(client MQTT.Client, message MQTT.Message) {
//	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
//}
