package files

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"

	"github.com/jodydadescott/shelly-client/sdk/types"
)

type File struct {
	STDIN    bool
	BaseName string
	FullName string
	Bytes    []byte
}

type Files struct {
	isDir bool
	Files []*File
}

func newFile(name string, bytes []byte, stdin bool) *File {
	f := &File{
		STDIN:    stdin,
		FullName: name,
		Bytes:    bytes,
	}

	if !stdin {
		f.BaseName = filepath.Base(name)
	}

	return f
}

func (t *Files) GetNamedFile() *File {

	if !t.isDir {
		return nil
	}

	for _, file := range t.Files {
		return file
	}

	return nil
}

func (t *Files) GetFile(name string) *File {

	getFile := func(subname string) *File {

		for _, file := range t.Files {
			if !file.STDIN {
				if file.BaseName == subname {
					return file
				}
				if file.BaseName == subname+".yaml" {
					return file
				}
				if file.BaseName == subname+".json" {
					return file
				}

				split := strings.Split(subname, "-")

				if file.BaseName == split[0]+".yaml" {
					return file
				}

				if file.BaseName == split[0]+".json" {
					return file
				}

			}
		}
		return nil
	}

	file := getFile(name)
	if file != nil {
		return file
	}

	file = getFile(strings.ToLower((name)))
	if file != nil {
		return file
	}

	file = getFile(strings.ToUpper((name)))
	if file != nil {
		return file
	}

	return nil
}

func (t *Files) GetSTDIN() *File {
	for _, file := range t.Files {
		if file.STDIN {
			return file
		}
	}
	return nil
}

func NewFiles(input string) (*Files, error) {

	files := &Files{}

	if input != "" {

		info, err := os.Stat(input)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {

			files.isDir = true

			rawFiles, err := os.ReadDir(input)
			if err != nil {
				return nil, err
			}

			for _, rawFile := range rawFiles {
				bytes, err := os.ReadFile(input + "/" + rawFile.Name())
				if err != nil {
					return nil, err
				}

				files.Files = append(files.Files, newFile(rawFile.Name(), bytes, false))
			}

		} else {

			bytes, err := os.ReadFile(input)
			if err != nil {
				return nil, err
			}

			files.Files = append(files.Files, newFile(input, bytes, false))

		}

		return files, nil
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {

		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		files.Files = append(files.Files, newFile("", bytes, true))

		return files, nil
	}

	return nil, fmt.Errorf("data input is required. Use filename or pipe to STDIN")
}

func (t *File) GetConfig() (*types.ShellyConfig, error) {

	var config types.ShellyConfig
	err := json.Unmarshal(t.Bytes, &config)
	if err == nil {
		return &config, nil
	}

	var errs *multierror.Error

	errs = multierror.Append(errs, err)

	err = yaml.Unmarshal(t.Bytes, &config)
	if err == nil {
		return &config, nil
	}

	errs = multierror.Append(errs, err)

	return nil, errs.ErrorOrNil()
}
