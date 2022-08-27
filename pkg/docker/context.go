package docker

import (
	"fmt"
	"os"
	"strings"

	"github.com/utkarsh-pro/kindli/pkg/sh"
)

// CreateContext will create a new docker context
func CreateContext(name, dockerHost string) error {
	return sh.RunSilent(fmt.Sprintf("docker context create %s --docker %s", name, dockerHost))
}

// DeleteContext deletes a docker context
func DeleteContext(name string) error {
	return sh.RunSilent(fmt.Sprintf("docker context delete %s", name))
}

// UseContext sets the given context as the default context
func UseContext(name string) error {
	os.Setenv("DOCKER_CONTEXT", name)
	return sh.RunSilent(fmt.Sprintf("docker context use %s", name))
}

// ExistsContext returns true if the given context already exists
func ExistsContext(name string) (bool, error) {
	resp, err := sh.RunIO("docker context ls -q")
	if err != nil {
		return false, err
	}

	for _, ctx := range strings.Split(string(resp), "\n") {
		if ctx == name {
			return true, nil
		}
	}

	return false, nil
}
