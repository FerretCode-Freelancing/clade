package registry

import (
	"os"
	"path/filepath"
)

func Remove(request Request, homeDir string) error {
	RegistryStore = Registry{}

	appDir := filepath.FromSlash(homeDir + "/.clade")

	if _, err := os.Stat(appDir); err != nil {
		return err
	}

	registryFile := filepath.FromSlash(appDir + "/registry.json")

	if _, err := os.Stat(registryFile); err != nil {
		return err
	}

	err := os.Remove(registryFile)
	if err != nil {
		return err
	}

	return nil
}
