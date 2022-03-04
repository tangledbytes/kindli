package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
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

	file, err := os.OpenFile(storePath, os.O_RDONLY|os.O_CREATE, 0660)
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
	if err := writeFS(); err != nil {
		logrus.Warn("failed to flush data to filesystem")
	}
}

// Get gets the value from the in memory store
func Get(key ...string) (interface{}, bool) {
	return utils.MapGet(_store, key...)
}

// DeleteTop deletes top level key from the map
func DeleteTop(key string) {
	delete(_store, key)
	if err := writeFS(); err != nil {
		logrus.Warn("failed to flush data to filesystem")
	}
}

func writeFS() error {
	file, err := os.Create(storePath)
	if err != nil {
		return fmt.Errorf("failed to write FS: %s", err)
	}

	defer file.Close()

	if err := json.NewEncoder(file).Encode(_store); err != nil {
		return fmt.Errorf("failed to write FS: %s", err)
	}

	return nil
}

// Flush writes the store data to the filesystem
func Flush() {
	utils.ExitIfNotNil(writeFS())
}
