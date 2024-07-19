package internal

import (
	"log/slog"
	"os"
)

const fName string = "\\__recent_cache__.json"

func GetCachePath() string {
	path, err := os.UserCacheDir()
	if err != nil {
		slog.Error("Can't get cache path: ", "error", err.Error())
		return ""
	}
	return path + fName
}
