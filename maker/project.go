package maker

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Project(fs fs.FS, name string) error {
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

	err = makeProject(fs, absPath, name)
	if err != nil {
		rerr := os.RemoveAll(absPath)
		if rerr != nil {
			return rerr
		}

		return err
	}

	return nil
}

func makeProject(fs fs.FS, absPath, name string) error {
	// check if folder already exists
	err := mainfile(fs, absPath, name)
	if err != nil {
		return err
	}

	err = defaultHTTPRoutes(fs, absPath)
	if err != nil {
		return err
	}

	err = defaultHomepage(fs, absPath)
	if err != nil {
		return err
	}

	err = defaultConfig(fs, absPath)
	if err != nil {
		return err
	}

	return nil

}

func mainfile(fs fs.FS, path, projectName string) error {
	path = filepath.Join(path, "cmd", projectName)

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(path, "main.go"))
	if err != nil {
		return err
	}

	tmpl, err := fs.Open("templates/main.go")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, tmpl)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// defaultHTTPRoutes will create a ./http dir with a http.go file.
// In this file we will register 1 route and 2 route handlers.
// The first handle is the home "/" handler.
// The second handler is a 404 handler.
func defaultHTTPRoutes(fs fs.FS, path string) error {
	path = filepath.Join(path, "internal", "http")

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(path, "http.go"))
	if err != nil {
		return err
	}

	tmpl, err := fs.Open("templates/http/http.go")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, tmpl)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func defaultHomepage(fs fs.FS, path string) error {
	err := os.Mkdir(filepath.Join(path, "template"), 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(path, "templates", "index.html"))
	if err != nil {
		return err
	}

	tmpl, err := fs.Open("templates/index.html")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, tmpl)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func defaultConfig(fs fs.FS, path string) error {
	f, err := os.Create(filepath.Join(path, "config.yml"))
	if err != nil {
		return err
	}

	tmpl, err := fs.Open("templates/config.yml")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, tmpl)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
