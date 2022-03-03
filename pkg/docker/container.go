package docker

import (
	"fmt"
	"strings"

	"github.com/utkarsh-pro/kindli/pkg/sh"
)

func ListRunningContainerNames() ([]string, error) {
	out, err := sh.RunIO("docker container ls --format={{.Names}}")
	if err != nil {
		return nil, fmt.Errorf("failed to get container list")
	}

	return strings.Split(string(out), "\n"), nil
}

func RunContainer(name string, rest string, removeExisting bool) error {
	if removeExisting {
		usable, err := IsContainerUsable(name)
		if err != nil {
			return err
		}

		if !usable {
			if err := RemoveNonUsableContainer(name); err != nil {
				return err
			}
		}
	}

	return sh.Run(fmt.Sprintf("docker run --name %s %s", name, rest))
}

func FindStoppedContainers() ([]string, error) {
	out, err := sh.RunIO("docker container ls -a -f 'status=exited' -f 'status=dead' -f 'status=created' -f 'status=paused' --format={{.Names}}")
	if err != nil {
		return nil, fmt.Errorf("failed to get container list")
	}

	return strings.Split(string(out), "\n"), nil
}

func IsContainerUsable(name string) (bool, error) {
	out, err := sh.RunIO(
		fmt.Sprintf(
			"docker container ls -a -f 'status=exited' -f 'status=dead' -f 'status=created' -f 'status=paused' -f 'name=%s' --format={{.Names}}",
			name,
		),
	)
	if err != nil {
		return false, fmt.Errorf("failed to find container")
	}

	return len(out) == 0, nil
}

func RemoveNonUsableContainer(name string) error {
	return sh.Run(fmt.Sprintf("docker container rm -f %s", name))
}
