package gwx

import (
	"embed"
	"io/fs"
	"os"
	"strings"
)

func WritePbToLocal(pb embed.FS) ([]string, error) {
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
			tmpFile, err := os.CreateTemp("", "*.pb")
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
