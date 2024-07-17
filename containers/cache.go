package containers

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func GetCache() ([]Request, error) {
	var requests []Request

	bytes, err := readCacheFile()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}

func WriteCache(toWrite []Request) error {
	requests, err := GetCache()
	if err != nil {
		return err
	}

	for i := range requests {
		for j := range toWrite {
			if requests[i].Name == toWrite[j].Name {
				toWrite = append(toWrite[:j], toWrite[:j+1]...)
			}
		}
	}

	requests = append(requests, toWrite...)

	bytes, err := json.Marshal(requests)
	if err != nil {
		return err
	}

	cacheFile, err := getCacheFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFile, bytes, os.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}

func getCacheFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.FromSlash(homeDir + "/.clade")

	if _, err := os.Stat(appDir); err != nil {
		return "", err
	}

	cacheFile := filepath.FromSlash(appDir + "/cache.json")

	return cacheFile, nil
}

func readCacheFile() ([]byte, error) {
	cacheFile, err := getCacheFilePath()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
