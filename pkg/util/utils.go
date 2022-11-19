package util

import "strings"

func LastPartFromPath(photoFilePath string) string {
	parts := strings.Split(photoFilePath, "/")
	return parts[len(parts)-1]
}
