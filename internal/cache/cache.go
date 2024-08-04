package cache

import (
	"bufio"
	"log/slog"
	"os"
	"sync"

	"github.com/Ha4sh-447/FiEx/internal"
	"github.com/Ha4sh-447/FiEx/pkg/files"
	"github.com/vmihailenco/msgpack/v5"
)

// TODO: Need to change the cache system
/*
 * Instead of storing search results in cache, cache the whole file system
 * as the application will search for files in the current dir
 * if the query is already present in cache it will show that result instead of the one found in current dir
 * hence cache the fs and then retrieve the dir from cache
 * and then search from the returned result
 */

type SearchCache struct {
	Store     map[string][]string `json:"store"`
	SyncStore *sync.Map
}

func NewSearchCache() *SearchCache {
	var s map[string][]string
	return &SearchCache{
		Store: s,
	}
}

func (r *SearchCache) WriteToFile_msgPack(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := msgpack.NewEncoder(writer)

	// Write data in chunks
	// for key, value := range r.Store {
	// 	if err := encoder.Encode(key); err != nil {
	// 		return err
	// 	}
	// 	if err := encoder.Encode(value); err != nil {
	// 		return err
	// 	}
	// }

	r.SyncStore.Range(func(key, value any) bool {
		if err := encoder.Encode(key); err != nil {
			return false
		}
		if err := encoder.Encode(value); err != nil {
			return false
		}
		return true
	})

	writer.Flush()
	return nil
}

func CreateSysCache() *SearchCache {
	sc := NewSearchCache()
	usr, err := os.UserHomeDir()
	if err != nil {
		slog.Error("ERROR", "Fetching Home Directory", err)
		return sc
	}

	f := files.TraverseDir(usr)
	if f != nil {
		sc.SyncStore = f
	}

	if err := sc.WriteToFile_msgPack(internal.GetCachePath()); err != nil {
		// if err := sc.WriteToFile_msgPack("output.msgpack"); err != nil {
		slog.Error("ERROR", "Writing to Cache File", err)
	} else {
		slog.Info("INFO", "SYS-CACHE", "Created Cache file")
	}
	return sc
}

func GetCache_msg(filename string) (*SearchCache, error) {
	file, err := os.Open(filename)
	if err != nil {
		slog.Error("ERROR", "OPENING FILE", err)
	}
	var cache SearchCache
	cache.SyncStore = &sync.Map{}
	reader := bufio.NewReader(file)
	decoder := msgpack.NewDecoder(reader)
	// cache.Store = make(map[string][]string)

	// Read data in chunks
	for {
		var key string
		var value []string
		if err := decoder.Decode(&key); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if err := decoder.Decode(&value); err != nil {
			return nil, err
		}
		// cache.Store[key] = value
		cache.SyncStore.Store(key, value)
	}

	return &cache, nil
}
