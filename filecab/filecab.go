package filecab

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Folder struct {
	STDIN    bool
	BaseName string
	FullName string
	Match    string
	Bytes    []byte
}

type Cabinet struct {
	FileOrDir string
	Match     []string
}

func NewCabinet(fileOrDir string) *Cabinet {
	return &Cabinet{
		FileOrDir: fileOrDir,
	}
}

func (t *Cabinet) SetFileOrDir(fileOrDir string) *Cabinet {
	t.FileOrDir = fileOrDir
	return t
}

func (t *Cabinet) AddMatch(match ...string) *Cabinet {
	t.Match = append(t.Match, match...)
	return t
}

func (t *Cabinet) IsDir() (bool, error) {
	info, err := os.Stat(t.FileOrDir)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func (t *Cabinet) Folder() (*Folder, error) {

	if t.FileOrDir == "" {

		fi, err := os.Stdin.Stat()
		if err != nil {
			return nil, err
		}

		if (fi.Mode() & os.ModeCharDevice) == 0 {

			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("file cabinet STDIN had error %w", err)
			}

			return newFolder("", bytes, true, ""), nil
		}

		return nil, fmt.Errorf("config must either be piped in from STDIN or the name of a file or directory must be specified")
	}

	info, err := os.Stat(t.FileOrDir)
	if err != nil {
		return nil, fmt.Errorf("file cabinet input file had error %w", err)
	}

	if !info.IsDir() {

		bytes, err := os.ReadFile(t.FileOrDir)
		if err != nil {
			return nil, fmt.Errorf("file cabinet read file had error %w", err)
		}

		return newFolder(t.FileOrDir, bytes, false, ""), nil
	}

	var folders []*Folder

	rawFiles, err := os.ReadDir(t.FileOrDir)
	if err != nil {
		return nil, err
	}

	for _, rawFile := range rawFiles {
		bytes, err := os.ReadFile(filepath.Join(t.FileOrDir + "/" + rawFile.Name()))
		if err != nil {
			return nil, err
		}

		folder := newFolder(rawFile.Name(), bytes, false, "")

		folders = append(folders, folder)
	}

	if len(folders) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("no file found in directory %s", t.FileOrDir))
	}

	if len(folders) == 1 {
		if len(t.Match) == 0 {
			folder := folders[0]
			folder.Match = "none specified"
			return folder, nil
		}
	}

	if len(t.Match) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("directory %s has %d files but no match pattern was specified", t.FileOrDir, len(t.Match)))
	}

	match := func(folder *Folder) bool {

		basename := strings.ToLower(folder.BaseName)

		for _, match := range t.Match {

			lowerMatch := strings.ToLower(folder.BaseName)

			if strings.Contains(basename, lowerMatch) {
				folder.Match = match
				return true
			}
		}

		return false
	}

	for _, folder := range folders {
		if match(folder) {
			return folder, nil
		}
	}

	return nil, fmt.Errorf("no matching file found")
}

func newFolder(name string, bytes []byte, stdin bool, match string) *Folder {

	if stdin {
		return &Folder{
			STDIN:    stdin,
			FullName: "STDIN",
			BaseName: "STDIN",
			Bytes:    bytes,
		}
	}

	return &Folder{
		FullName: name,
		BaseName: filepath.Base(name),
		Bytes:    bytes,
		Match:    match,
	}
}
