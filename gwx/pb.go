package gwx

import (
	"embed"
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
