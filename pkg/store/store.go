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
	_store   map[string]interface{} = nil
	_idIndex []bool                 = make([]bool, 100)

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

	for _, ent := range _store {
		castedEnt, ok := ent.(map[string]interface{})
		if !ok {
			logrus.Debug("failed to cast store entity")
			continue
		}

		instanceID, ok := castedEnt["instanceID"]
		if !ok {
			logrus.Debug("failed to get instance ID")
			continue
		}

		castedInstanceID, ok := instanceID.(float64)
		if !ok {
			logrus.Debug("failed to cast extracted instance id")
			continue
		}

		_idIndex[int(castedInstanceID)] = true
		logrus.Debug("Found instance id: ", int(castedInstanceID))
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

// GetNextID returns the next ID that can be used for the next instance
func GetNextID() (int, error) {
	for i := 0; i < len(_idIndex); i++ {
		if !_idIndex[i] {
			return i, nil
		}
	}

	return -1, fmt.Errorf("all the ids are taken")
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
