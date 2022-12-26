package envHelp

import "os"

func BaseDirFromEnv() string {
	var basedir string
	if os.Getenv("I_AM_EMBEDDED") == "true" {
		basedir = ""
	} else {
		basedir = "."
	}
	return basedir
	//
	//environment := os.Getenv("ENVIRONMENT")
	//var baseDir string
	//if environment == "dev" {
	//	baseDir = "."
	//} else {
	//	baseDir = ""
	//}
	//return baseDir
}

func NodeNameFromEnv() string {
	nodeName := os.Getenv("MY_NODE_NAME")
	if nodeName == "" {
		nodeName = "unspecifiedNodeName"
	}
	return nodeName
}
