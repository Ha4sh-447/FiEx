package cache

import (
	"encoding/json"
	"os"
)

type SearchCache struct {
	Store map[string][]string `json:"store"`
}

func NewSearchCache() *SearchCache {
	var s map[string][]string
	return &SearchCache{
		Store: s,
	}
}

func (r *SearchCache) Add(query string, res []string) {
	if _, exists := r.Store[query]; !exists {
		r.Store[query] = res
	}
	// }
}

// loads the entire cache
func GetCache(cachePath string) (*SearchCache, error) {
	file, err := os.Open(cachePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() == 0 {
		// Return a SearchCache with an empty map if the file is empty
		return &SearchCache{Store: make(map[string][]string)}, nil
	}

	var data SearchCache
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data.Store); err != nil {
		return nil, err
	}
	// fmt.Println(reflect.TypeOf(data))
	return &data, nil
}

func (r *SearchCache) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// slog.Info("Writing to recent cache: ", r.Store)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(r.Store); err != nil {
		return err
	}

	return nil
}
