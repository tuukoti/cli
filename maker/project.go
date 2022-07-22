package maker

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func Project(fs fs.FS, pkgURL string) error {
	name := filepath.Base(pkgURL)

	_, err := os.Stat(name)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if err == nil {
		return fmt.Errorf("project folder already exists")
	}

	absPath, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	err = os.Mkdir(absPath, 0777)
	if err != nil {
		return fmt.Errorf("failed to create the project folder, %w", err)
	}

	err = makeProject(fs, absPath, pkgURL, name)
	if err != nil {
		rerr := os.RemoveAll(absPath)
		if rerr != nil {
			return rerr
		}

		return err
	}

	return nil
}

func makeProject(fs fs.FS, absPath, pkgURL, name string) error {
	// main file
	err := copyFile(
		fs,
		"templates/main.go.tmpl",
		filepath.Join(absPath, "cmd", name, "main.go"),
		func(t []byte) ([]byte, error) {
			return bytes.ReplaceAll(t, []byte("{{resources_pkg}}"), []byte(pkgURL+"/resources")), nil
		},
	)
	if err != nil {
		return err
	}

	// resources
	err = copyFile(
		fs,
		"templates/resources/resources.go",
		filepath.Join(absPath, "resources", "resources.go"),
		nil,
	)
	if err != nil {
		return err
	}

	err = copyFile(
		fs,
		"templates/resources/default.go",
		filepath.Join(absPath, "resources", "default.go"),
		nil,
	)
	if err != nil {
		return err
	}

	// config file
	err = copyFile(
		fs,
		"templates/config.yml",
		filepath.Join(absPath, "config.yml"),
		nil,
	)
	if err != nil {
		return err
	}

	// default html files
	err = copyFile(
		fs,
		"templates/views/index.html",
		filepath.Join(absPath, "views", "index.html"),
		nil,
	)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "init", pkgURL)
	cmd.Dir = absPath
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = absPath
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil

}

type modifierFunc func(t []byte) ([]byte, error)

func copyFile(fs fs.FS, src, dst string, modifier modifierFunc) error {
	tmpl, err := fs.Open(src)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(tmpl)
	if err != nil {
		return err
	}

	b := buf.Bytes()

	if modifier != nil {
		b, err = modifier(b)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(filepath.Dir(dst), 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil

}
