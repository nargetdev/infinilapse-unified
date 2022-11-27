package compiler

import (
	"fmt"
	"github.com/bitfield/script"
	"infinilapse-unified/pkg/envHelp"
	"log"
	"os"
)

func CompileAllPreviousVideo() (string, error) {
	basedir := envHelp.BaseDirFromEnv()
	sourceDir := basedir + "/data/out"

	var chunkSourceDirectoriesList []string = streamSourcePaths(sourceDir)
	fmt.Printf("chunkSourceDirectoriesList []string --- %v\n", chunkSourceDirectoriesList)

	for _, chunkDir := range chunkSourceDirectoriesList {
		println(chunkDir)
		chunksSlice := chunksFromPath(chunkDir)
		//fmt.Printf("chunksSlice []string --- %v\n", chunksSlice)

		outMp4Path := chunkDir + ".latest.mp4"

		chunksOutDir := basedir + "/data/meta/" + lastPartFromPath(chunkDir)
		chunksSliceTxtPath := chunksOutDir + "/chunksSlice.txt"
		fmt.Printf("About to write %s\n", chunksSliceTxtPath)
		_ = constructChunksSliceTxtFromChunksSlice(chunksSlice, chunksOutDir)

		chunksSliceTxtString, readErr := os.ReadFile(chunksSliceTxtPath)
		check(readErr, "readErr")

		fmt.Printf("=== chunksSlice.txt ===\n%s\n", chunksSliceTxtString)

		var inputChunksOpts string = ""
		for _, aChunkPath := range chunksSlice {
			inputChunksOpts += " -i " + aChunkPath
		}

		var stitcherCmd string = fmt.Sprintf("ffmpeg -y -f concat -safe 0 -i "+chunksSliceTxtPath+" -c copy %s", outMp4Path)
		//var stitcherCmd string = fmt.Sprintf("ffmpeg -f concat -safe 0 "+inputChunksOpts+" -c copy %s", outMp4Path)
		stitchStdout, stitchErr := script.Exec(stitcherCmd).String()
		if DEBUG {
			println(stitcherCmd)
			fmt.Printf("%s\n", stitchStdout)
		}
		check(stitchErr, "stitchErr")

	}

	return "success", nil
}

func check(e error, label string) {
	if e != nil {
		fmt.Printf("check({{ %s }}) --- %s\n", label, e)
	}
}

func constructChunksSliceTxtFromChunksSlice(listOfChunkPath []string, outDir string) (textWritten string) {
	// start by ensuring directory is available
	checkStdoutErr(script.Exec("mkdir -p " + outDir))

	outStr := ""
	for _, chunkPath := range listOfChunkPath {
		outputLine := fmt.Sprintf("file %s", chunkPath)
		outStr = fmt.Sprintf("%s\n", outStr+outputLine)
	}
	err := os.WriteFile(outDir+"/chunksSlice.txt", []byte(outStr), 0644)
	check(err, "os.WriteFile")

	return outStr
}

func checkStdoutErr(exec *script.Pipe) {
	_, err := exec.Stdout()
	if err != nil {
		fmt.Printf("::ERR:: %s", err)
		return
	}
}

func chunksFromPath(chunkDir string) (chunks []string) {
	files, err := os.ReadDir(chunkDir)
	if err != nil {
		log.Printf("%s\n", err)
	}
	for _, file := range files {
		if file.IsDir() {
			_ = fmt.Errorf("this should have been an mp4 file")
		} else if file.Name() == ".DS_Store" {
			_ = fmt.Errorf("ignore .DS_Store")
		} else {
			chunks = append(chunks, chunkDir+"/"+file.Name())
		}
	}
	return chunks
}

func streamSourcePaths(chunkDir string) (streamDevice []string) {
	files, err := os.ReadDir(chunkDir)
	if err != nil {
		log.Printf("%s\n", err)
	}

	for _, file := range files {
		if file.IsDir() {
			streamDevice = append(streamDevice, chunkDir+"/"+file.Name())
		}
	}

	return
}
