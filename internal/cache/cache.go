package cache

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/Ha4sh-447/FiEx/internal"
	"github.com/Ha4sh-447/FiEx/pkg/files"
	"github.com/vmihailenco/msgpack/v5"
)

// TODO: Need to change the cache system
/*
* If any changes occur to the file system,
* read the changes and update the cache accoding to it
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

// This caches just the user's home directory,
// this is because it takes long time to cache it
// hence, would be faster to just search the ones which aren't
// present in the cache and later add them
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

func (sc *SearchCache) Add(key string, val []string) {
	sc.SyncStore.Store(key, val)
	slog.Info("Added to cache", "Success", fmt.Sprintf("Added: %s", key))
}

func (sc *SearchCache) Update(key string, val []string) {
	if res, present := sc.SyncStore.Load(key); present {
		for _, v := range val {
			sc.SyncStore.Store(key, append(res.([]string), v))
		}
	} else {
		slog.Error("Key not present", "Add key", key)
	}
}
