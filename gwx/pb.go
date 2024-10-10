package gwx

import (
	"embed"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"strings"
)

type Opts func(config *gwxConfig)

type gwxConfig struct {
	Dir string
}

func WritePbToLocal(pb embed.FS, opts ...Opts) ([]string, error) {
	config := &gwxConfig{}

	for _, opt := range opts {
		opt(config)
	}

	var fileList []string

	err := fs.WalkDir(pb, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".pb") {
			data, err := pb.ReadFile(path)
			if err != nil {
				return err
			}
			if stat, err := os.Stat(config.Dir); err != nil {
				if !os.IsExist(err) {
					err = os.MkdirAll(config.Dir, 0o755)
					if err != nil {
						return err
					}
				}
			} else {
				if !stat.IsDir() {
					return errors.Errorf("%s: not a directory", config.Dir)
				}
			}

			tmpFile, err := os.CreateTemp(config.Dir, "*.pb")
			if err != nil {
				return err
			}
			defer tmpFile.Close()
			if _, err := tmpFile.Write(data); err != nil {
				return err
			}
			fileList = append(fileList, tmpFile.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}
