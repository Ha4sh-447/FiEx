package internal

import (
	"log/slog"
	"os"
)

const fName string = "\\__recent_cache__.msgpack"

func GetCachePath() string {
	path, err := os.UserCacheDir()
	if err != nil {
		slog.Error("Can't get cache path: ", "error", err.Error())
		return ""
	}
	return path + fName
}

func MatchString(s, m string) int {
	return -1
}
