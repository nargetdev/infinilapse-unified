package envHelp

import "os"

func BaseDirFromEnv() string {
	environment := os.Getenv("ENVIRONMENT")
	var baseDir string
	if environment == "dev" {
		baseDir = "."
	} else {
		baseDir = ""
	}
	return baseDir
}
