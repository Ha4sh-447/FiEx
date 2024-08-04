package pkg

import (
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	internal "github.com/Ha4sh-447/FiEx/internal/cache"
)

type topRes struct {
	Path  string
	Score int
}

func Search(dir, query string, res []string) []string {
	start := time.Now()

	// fileBuff, err := TraverseDir(dir)
	// if err != nil {
	// 	log.Fatal("Can't read file", err)
	// }

	var topResult []topRes
	var wg sync.WaitGroup
	results := make(chan topRes, 50)

	maxSets := 50
	dataPerSet := len(res) / maxSets
	fmt.Println("Len and dps: ", len(res), dataPerSet)
	in := 0
	end := dataPerSet
	// slog.Info("CurrRead: ", currRead)
	for i := 0; i < maxSets; i++ {
		data := res[in:end]
		wg.Add(1)
		go func(data []string) {
			defer wg.Done()
			sc := maxScore(data, query)
			results <- sc
		}(data)
		end += dataPerSet
		in += dataPerSet
	}

	wg.Wait()
	close(results)
	fmt.Println(len(results))

	for r := range results {
		if r.Score != 0 {
			topResult = append(topResult, r)
		}
	}

	fmt.Println(len(topResult))

	sort.Slice(topResult, func(i, j int) bool {
		return topResult[i].Score > topResult[j].Score
	})

	// topResults := make([]string, 50)
	var topResults []string

	for _, i := range topResult {
		topResults = append(topResults, i.Path)
	}

	slog.Info("Time taken to complete search: ", "Info", time.Since(start))
	return topResults
}

func maxScore(fileBuff []string, query string) topRes {

	m := math.MinInt
	var str string
	for _, l := range fileBuff {
		t := m
		score := Score(l, query)
		m = int(math.Max(float64(score), float64(m)))
		if m != t {
			str = l
		}
	}
	str = strings.ReplaceAll(str, "\\\\", "\\")

	return topRes{Path: str, Score: m}
}

// TODO: Better scoring system required
/*
	* Exact match => 45
	* Match at end and start - exact matches, word to word => 30
	* At the end, if the end part of the string if there is not an exact match but still a match => 25
	* For searching as prefix =>
			* Start string checking relative to current directory path
			* scoring system = (20, 15)
	* For matching within string =>
			* get the maximum subsequence that matches the query
			* score -> length * 2
*/

func Score(path, query string) int {
	var score int = 0
	path_l := strings.ToLower(path)
	query_l := strings.ToLower(query)

	// pL := len(path_l)

	store := strings.Split(path_l, "\\")
	// a way to calculate difference in string paths
	// get that part of the query which is after the path string

	// relPath := strings.Split(path_l, query_l)

	// fmt.Println("Path: ", path_l)
	// fmt.Println("Query: ", query_l)
	// fmt.Println("Path ending: ", store[len(store)-1])
	end := store[len(store)-1]
	// fmt.Println("Split the path after the query: ", relPath)

	// Match at the end
	if strings.HasSuffix(end, query_l) {
		score += 52
	} else {
		score += 50
	}

	// Match at the beginning
	if strings.HasPrefix(path, query) {
		score += 22
	} else {
		score += 20
	}

	// Match in the middle
	if !strings.HasPrefix(path, query) && !strings.HasSuffix(path, query) && (strings.Contains(path, query) || strings.Contains(path_l, query_l)) {
		score += 10
	}

	// Match a substring within the string
	r1, r2 := []rune(path), []rune(query)
	count := 0
	maxCount := 0

	j := 0
	for i := range r1 {
		if j < len(r2) && r1[i] == r2[j] {
			count++
			j++
		} else {
			if count > maxCount {
				maxCount = count
			}
			count = 0
			j = 0
		}
	}
	if count > maxCount {
		maxCount = count
	}
	if maxCount == 1 {
		maxCount = 0
	}
	score += (5 * maxCount)

	return score
}

func SearchInCache(currDir string, res *internal.SearchCache) []string {

	if res.SyncStore == nil {
		return nil
	}
	var result []string

	res.SyncStore.Range(func(key, value any) bool {
		if key == currDir {
			result = value.([]string)
			slog.Info("Found in Cache", "dir", key)
			return false
		}
		return true
	})

	// for q, r := range res.Store {
	// 	if q == filepath.Clean(currDir) {
	// 		slog.Info("Found in Cache", "dir", q)
	// 		fmt.Println(r)
	// 		return r
	// 	}
	// }

	slog.Info("Not in cache", "dir", currDir)
	return result
}
