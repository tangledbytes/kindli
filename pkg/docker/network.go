package docker

import (
	"fmt"
	"strings"

	"github.com/utkarsh-pro/kindli/pkg/sh"
)

// NetworkInsepct will inspect a docker network and will return the response as per the format string
func NetworkInspect(network, format string) (string, error) {
	resp, err := sh.RunIO(fmt.Sprintf("docker network inspect %s -f='%s'", network, format))
	if err != nil {
		return "", fmt.Errorf("failed to inspect docker network: %s", err)
	}

	return strings.Trim(string(resp), " \n"), nil
}
