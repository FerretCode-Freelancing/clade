package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Add(request Request, homeDir string) error {
	RegistryStore = Registry{
		Username: request.Registry.Username,
		Secret:   request.Registry.Secret,
	}

	if !request.Store {
		return nil
	}

	// TODO: inform user the secret will be stored in plaintext

	appDir := filepath.FromSlash(homeDir + "/.clade")

	if _, err := os.Stat(appDir); err != nil {
		return err
	}

	data, err := json.Marshal(request.Registry)
	if err != nil {
		return err
	}

	err = os.WriteFile(
		filepath.FromSlash(appDir+"/registry.json"),
		data,
		os.ModeAppend,
	)
	if err != nil {
		return err
	}

	return nil
}
