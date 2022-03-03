package registry

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/docker"
)

//go:embed registry-cache-config.template
var cfg string

type Registry struct {
	Name string
	URL  string
	Port string
}

func Knowns() []Registry {
	return []Registry{
		{
			Name: "dockerio-registry",
			URL:  "https://registry-1.docker.io",
			Port: "5000",
		},
		{
			Name: "quayio-registry",
			URL:  "https://quay.io",
			Port: "5010",
		},
		{
			Name: "gcrio-registry",
			URL:  "https://gcr.io",
			Port: "5020",
		},
	}
}

func New(name, url, port string) *Registry {
	return &Registry{
		Name: name,
		URL:  url,
		Port: port,
	}
}

func (r *Registry) Create(cacheDir string) error {
	cfgPath, err := createConfig(r.Name, map[string]interface{}{
		"RemoteURL": r.URL,
		"Port":      r.Port,
	})
	if err != nil {
		return err
	}

	return docker.RunContainer(
		r.Name,
		fmt.Sprintf("-d --restart=always -v %s:/etc/docker/registry/config.yaml -v %s:/var/lib/registry registry:2", cfgPath, cacheDir),
		true,
	)
}

func (r *Registry) IsRunning() (bool, error) {
	runningContainers, err := docker.ListRunningContainerNames()
	if err != nil {
		return false, err
	}

	for _, cnt := range runningContainers {
		if cnt == r.Name {
			return true, nil
		}
	}

	return false, nil
}

func createConfig(name string, overrides map[string]interface{}) (string, error) {
	path := filepath.Join(config.Dir(), fmt.Sprintf("%s.yaml", name))

	parsed, err := template.New(name).Parse(cfg)
	if err != nil {
		return path, err
	}

	file, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer file.Close()

	return path, parsed.Execute(file, overrides)
}
