package cmd

import (
	"os"
	"path"
)

func getConfigFileLocation() (string, bool) {
	home := getOrPanic(os.UserHomeDir())

	dir := path.Join(home, ".aptos")

	if _, err := os.Stat(path.Join(dir, "config.yaml")); !os.IsNotExist(err) {
		return path.Join(dir, "config.yaml"), true
	}
	if _, err := os.Stat(path.Join(dir, "config.yml")); !os.IsNotExist(err) {
		return path.Join(dir, "config.yml"), true
	}

	return path.Join(dir, "config.yaml"), false
}
