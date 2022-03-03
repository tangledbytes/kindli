package docker

import (
	"fmt"

	"github.com/utkarsh-pro/kindli/pkg/sh"
)

// Connect a network with a container
func NetworkConnect(network, container string) error {
	return sh.Run(fmt.Sprintf("docker network connect %s %s", network, container))
}
