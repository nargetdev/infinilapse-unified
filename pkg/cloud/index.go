package cloud

import (
	"fmt"
	"infinilapse-unified/pkg/gcpMgmt"
	"infinilapse-unified/pkg/gqlMgmt"
	"os"
	"strings"
)

func IndexGoogleCloudStorageAndGraphQL(files []string) error {
	var nodeName string = nodeNameFromEnv()
	var bucket string = bucketFromEnv()
	for _, file := range files {

		//bucket := "tl-data"
		objPath := nodeName + "/" + file

		objUrl, _ := gcpMgmt.StoreFileToBucket(file, objPath, bucket)

		// upload was success index in the graph
		gqlMgmt.IndexToGraph(objUrl, bucket, nodeName+"."+cameraNameFromFullPath(file))
	}

	return nil
}

func cameraNameFromFullPath(file string) string {
	parts := strings.Split(file, "/")
	return parts[len(parts)-2]
}

func bucketFromEnv() string {
	GcbBucket := "gcb-site-pub"
	return GcbBucket
}

func nodeNameFromEnv() string {
	nodeName := os.Getenv("MY_NODE_NAME")
	if nodeName == "" {
		nodeName = "unspecifiedNodeName"
	}
	return nodeName
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
