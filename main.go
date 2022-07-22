package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/tuukoti/cli/maker"
)

const colorRed = "\033[0;31m"
const colorGreen = "\033[0;32m"
const colorBlue = "\033[0;34m"
const colorYellow = "\033[93m"
const colorNone = "\033[0m"
const newLine = "\n"

var Version = "develop"

//go:embed templates
var fs embed.FS

func main() {
	if err := runCLI(); err != nil {
		fmt.Printf("%s%s%s", colorRed, err, newLine)
	}

	fmt.Println("")
}

func runCLI() error {

	fmt.Println("------------------")
	fmt.Printf("Tuukoti CLI - %s%s", Version, newLine)
	fmt.Println("------------------")
	fmt.Println("")

	args := os.Args
	if len(args) == 1 {
		err := fmt.Errorf("missing command, try `tuukoti list` to see what you can run.")

		fmt.Println(err)

		return nil
	}

	switch os.Args[1] {
	case "list":
		printCmd("new sitename", "create a new project within the directory of sitename")
		printCmd("serve", "serves your app locally")
		printCmd("build", "build the binary and any resources configured")
		printCmd("list", "list displays a list of all commands")
		printCmd("version", "displays the cli tools current version")
		fmt.Println("")
		printHeader("make")
		printCmd("make resource", "creates a new route file and a ")
	case "version":
		fmt.Println("")
		version()
	case "new":
		err := maker.Project(fs, os.Args[2])
		if err != nil {
			return err
		}

		fmt.Printf("%sproject %s has been created%s", colorGreen, os.Args[2], newLine)
	case "serve":

	case "build":

	case "make":
		if len(os.Args) == 2 {
			fmt.Printf("missing what we should make %s", newLine)
			return nil
		}

		switch os.Args[2] {
		case "resource":
			resourceName := os.Args[3]

			err := maker.Resource(resourceName)
			if err != nil {
				return err
			}
		}

	default:
		err := fmt.Errorf("missing command, try `tuukoti list` to see what you can run.")

		fmt.Println(err)
	}

	return nil
}

func version() {
	fmt.Printf("Version: %s%s", Version, newLine)
}

func printCmd(commandName, desc string) {
	fmt.Printf("  %s%-12s %s- %s.%s", colorGreen, commandName, colorNone, desc, newLine)
}

func printHeader(h string) {
	fmt.Printf("%s%s%s", colorYellow, h, newLine)
}
