package config

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/knadh/koanf/parsers/yaml"
)

const configDir = "../../configs"

func TestConfigs(t *testing.T) {
	configFS := os.DirFS(configDir)
	err := fs.WalkDir(configFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		cFile, err := configFS.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", path, err)
		}

		content, err := io.ReadAll(cFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		_, err = yaml.Parser().Unmarshal(content)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestKoanfLoad(t *testing.T) {
}
