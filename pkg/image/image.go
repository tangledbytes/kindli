package image

import (
	"fmt"

	"github.com/utkarsh-pro/kindli/pkg/sh"
)

func Load(image, cluster string) error {
	if cluster == "" {
		cluster = "kindli"
	}

	if err := sh.Run(fmt.Sprintf("kind load docker-image %s --name %s", image, cluster)); err != nil {
		return fmt.Errorf("failed to load image into the cluster \"%s\": %s", cluster, err)
	}

	return nil
}
