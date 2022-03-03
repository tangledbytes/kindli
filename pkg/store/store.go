package store

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var (
	_store map[string]interface{} = nil

	storeFile = "state.json"
	storePath = ""
)

// Load the store file into the memory
func Load() {
	// noop if the store has been initialized once
	if _store != nil {
		return
	}

	storePath = filepath.Join(config.Dir(), storeFile)

	file, err := os.OpenFile(storePath, os.O_RDWR|os.O_CREATE, 0660)
	utils.ExitIfNotNil(err)
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&_store); err != nil {
		if err == io.EOF {
			_store = map[string]interface{}{}
			return
		}

		utils.ExitIfNotNil(err)
	}
}

// Set sets value to in memory store
func Set(value interface{}, key ...string) {
	utils.MapSet(_store, value, key...)
}

// Get gets the value from the in memory store
func Get(key ...string) (interface{}, bool) {
	return utils.MapGet(_store, key...)
}

// DeleteTop deletes top level key from the map
func DeleteTop(key string) {
	delete(_store, key)
}

// Flush writes the store data to the filesystem
func Flush() {
	file, err := os.OpenFile(storePath, os.O_RDWR|os.O_CREATE, 0660)
	utils.ExitIfNotNil(err)
	defer file.Close()

	utils.ExitIfNotNil(json.NewEncoder(file).Encode(_store))
}
