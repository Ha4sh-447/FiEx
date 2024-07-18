package files

import (
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	internal "github.com/Ha4sh-447/FiEx/internal/cache"
	"github.com/saracen/walker"
)

func walkerReadFiles(path string) ([]string, error) {
	fmt.Println("Starting file search")
	start := time.Now()
	var fileBuff []string

	walkFun := func(pathname string, fi os.FileInfo) error {
		fileBuff = append(fileBuff, pathname)
		return nil
	}

	// cpuLimit := walker.WithLimit(25)

	// error callback option
	errorCallbackOption := walker.WithErrorCallback(func(pathname string, err error) error {
		// ignore permissione errors
		if os.IsPermission(err) {
			return nil
		}
		// halt traversal on any other error
		return err
	})

	walker.Walk(path, walkFun, errorCallbackOption)
	slog.Info("Time Taken by Walker function: ", time.Since(start))
	return fileBuff, nil
}

type topRes struct {
	Path  string
	Score int
}

func Search(dir, query string) []string {
	start := time.Now()

	fileBuff, err := walkerReadFiles(dir)
	if err != nil {
		log.Fatal("Can't read file", err)
	}

	var topResult []topRes
	var wg sync.WaitGroup
	results := make(chan topRes, 50)

	maxSets := 50
	dataPerSet := len(fileBuff) / maxSets
	fmt.Println("Len and dps: ", len(fileBuff), dataPerSet)
	in := 0
	end := dataPerSet
	// slog.Info("CurrRead: ", currRead)
	for i := 0; i < maxSets; i++ {
		data := fileBuff[in:end]
		wg.Add(1)
		go func(data []string) {
			defer wg.Done()
			sc := maxScore(data, query)
			results <- sc
		}(data)
		end += dataPerSet
		in += dataPerSet
		// fmt.Println("Iterations res data len: ", len(data))
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
	fmt.Println(topResult)

	// topResults := make([]string, 50)
	var topResults []string

	for _, i := range topResult {
		// if i.Path != "" {
		topResults = append(topResults, i.Path)
		// }
	}

	slog.Info("Time taken to complete search: ", time.Since(start))
	// fmt.Println(topResult)
	fmt.Println(len(topResults))
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
	fmt.Println("High score string: ", str)
	str = strings.ReplaceAll(str, "\\\\", "\\")
	// slog.Info("High score string: ", str)
	return topRes{Path: str, Score: m}
}

// Better scoring system required
func Score(path, query string) int {
	var score int = 0
	path_l := strings.ToLower(path)
	query_l := strings.ToLower(query)

	pL := len(path_l)

	// Match at the end
	if strings.HasSuffix(path_l, query_l) {
		if strings.HasSuffix(path_l, query_l) {
			score += 52
		} else {
			score += 50
		}
	}

	// Match at the beginning
	if strings.HasPrefix(path, query) {
		if strings.HasSuffix(path_l, query_l) {
			score += 32
		} else {
			score += 30
		}
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

	return score - pL
}

func SearchRecentcache(query, path string) []string {
	res, err := internal.GetCache(path)
	if err != nil {
		slog.Error("Cache error: ", err)
		return nil
	}

	for _, r := range res {
		if p := r[query]; p != nil {
			return p
		}
	}

	slog.Info(query, " not in cache")
	return nil
}
