package main

import (
	"fmt"
	"infinilapse-unified/pkg/compiler"
)

func main() {
	// stitch together all previous.
	var result string
	var stitchErr error
	result, stitchErr = compiler.CompileAllPreviousVideo()
	if stitchErr != nil {
		fmt.Printf(":::ERR:::  stitch.CompileAllPreviousVideo()\n%s\n\n:::RESULT:::\n%s\n\n", stitchErr, result)
	} else {
		fmt.Printf(":::RESULT:::  stitch.CompileAllPreviousVideo()\n%s\n\n", result)
	}

}
