package utl

import "strings"

func TrimPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return path[1:]
	}
	return path
}
