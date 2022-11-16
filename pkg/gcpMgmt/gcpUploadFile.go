package gcpMgmt

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// UploadFile uploads an object.
func UploadFile(filePath, object, bucket string) error {
	fmt.Printf("Uploading filePath.. %s\tto bucket.. %s at obj path... %s\n", filePath, bucket, object)

	//fileName := "getToDaBucket/1666100095.jpg"
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v\n", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	//// do latest
	//latest := client.Bucket(bucket).Object("latest.jpg")
	//// Upload an object with storage.Writer.
	//latestWc := latest.NewWriter(ctx)
	//if _, err = io.Copy(latestWc, f); err != nil {
	//	return fmt.Errorf("io.Copy: %v", err)
	//}
	//if err := latestWc.Close(); err != nil {
	//	return fmt.Errorf("Writer.Close: %v", err)
	//}
	////fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	//fmt.Printf("LATEST Blob %v uploaded.\n", object)

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	//o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	//attrs, err := o.Attrs(ctx)
	//if err != nil {
	//     return fmt.Errorf("object.Attrs: %v", err)
	//}
	//o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	//fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	//fmt.Printf("Blob %v uploaded.\n", object)

	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %v", err)
	}

	fmt.Printf("Blob %v is now publicly accessible.\n", object)

	return nil
}

func FileIsEmpty(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("%s", err)
		time.Sleep(1)
		return true
	}
	fi, err := file.Stat()
	if err != nil {
		fmt.Printf("%s", err)
		time.Sleep(1)
		return true
	}
	miSize := fi.Size()
	fmt.Println(miSize)
	return miSize == 0
}
