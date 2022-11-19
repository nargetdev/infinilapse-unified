package gcpMgmt

import (
	"fmt"
)

func StoreFileToBucket(fullPath, objPath, bucket string) (objUrl string, err error) {

	// hai o cool
	//UploadFile
	errr := UploadFile(fullPath, objPath, bucket)
	if errr != nil {
		fmt.Printf("\n\nERROR UPLOADING FILE:")
		fmt.Println(errr)
		return "error", errr
	}

	objUrl = "https://storage.googleapis.com/" + bucket + "/" + objPath
	println("\n\nOBJECT URL:\n\n" + objUrl)

	return objUrl, nil
}
