package files

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// IMPLEMENTING A VERY BASIC FUZZY SEARCH ALGORITHM
// USING LEVENSHTEIN DISTANCE

// TODO: CHANGE TO CUSTOM IMPELEMENTATION FOR BETTER SCORING SYSTEM
func LevenshteinDistance(s1, s2 string) int {
	lenS1 := len(s1)
	lenS2 := len(s2)

	if lenS1 < lenS2 {
		return LevenshteinDistance(s2, s1)
	}

	previousRow := make([]int, lenS2+1)
	for i := range previousRow {
		previousRow[i] = i
	}

	for i := 1; i <= lenS1; i++ {
		currentRow := make([]int, lenS2+1)
		currentRow[0] = i
		for j := 1; j <= lenS2; j++ {
			insertions := previousRow[j] + 1
			deletions := currentRow[j-1] + 1
			substitutions := previousRow[j-1]
			if s1[i-1] != s2[j-1] {
				substitutions++
			}
			currentRow[j] = min(insertions, min(deletions, substitutions))
		}
		previousRow = currentRow
	}

	return previousRow[lenS2]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ExpandHomeDir expands the ~ to the home directory
func ExpandHomeDir(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// GetAllFiles recursively collects all file paths starting from the given directory.
func GetAllFiles(directory string) ([]string, error) {
	var filePaths []string
	err := filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	return filePaths, err
}

// FuzzySearchFiles searches for files matching the query within the given directory.
func FuzzySearchFiles(directory, searchQuery string) ([]string, error) {
	filePaths, err := GetAllFiles(directory)
	if err != nil {
		return nil, err
	}

	type match struct {
		distance int
		path     string
	}
	var matches []match

	for _, filePath := range filePaths {
		fileName := filepath.Base(filePath)
		distance := LevenshteinDistance(fileName, searchQuery)
		matches = append(matches, match{distance, filePath})
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].distance < matches[j].distance
	})

	var sortedMatches []string
	for _, m := range matches {
		sortedMatches = append(sortedMatches, m.path)
	}

	return sortedMatches, nil
}
