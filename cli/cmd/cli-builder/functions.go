package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alioth-center/infrastructure/cli"
	"github.com/alioth-center/infrastructure/utils/values"
)

// InitProject is a handler for the "init" command
// command: init <directory> <package-name>
func InitProject(input *cli.Input) {
	if !checkInitProject(input) {
		return
	}
}

func LoadProject(input *cli.Input) {
}

func checkInitProject(input *cli.Input) bool {
	directory := input.Params.GetString("directory")
	packageName := input.Params.GetString("package-name")
	if directory == "" || packageName == "" {
		_, text := i18nCollections[i18nInitProjectBadParams].GetTranslation(input.Language...)
		fmt.Println(text)
		return false
	}

	args := map[string]string{
		"directory": directory,
	}

	// check if directory exists and is a directory
	f, readDir := os.Stat(directory)
	if readDir != nil && !os.IsNotExist(readDir) {
		_, text := i18nCollections[i18nInitProjectReadWorkingDirectoryError].GetTranslation(input.Language...)
		args["error"] = readDir.Error()
		fmt.Println(values.NewStringTemplate(text, args).Parse())
		return false
	}
	if readDir != nil && os.IsNotExist(readDir) {
		// not exist, try to create the directory and rejudge
		_, text := i18nCollections[i18nInitProjectCreateWorkingDirectoryCheck].GetTranslation(input.Language...)
		fmt.Print(values.NewStringTemplate(text, args).Parse())
		var inputWord string
		scanErr := errors.New("nil err")
		for scanErr != nil {
			_, scanErr = fmt.Scanln(&inputWord)
		}
		if !strings.HasPrefix(strings.ToLower(inputWord), "y") {
			_, text := i18nCollections[i18nInitProjectCreateWorkingDirectoryAbort].GetTranslation(input.Language...)
			fmt.Println(text)
			return false
		}

		createErr := os.MkdirAll(directory, 0o755)
		if createErr != nil {
			_, text := i18nCollections[i18nInitProjectCreateWorkingDirectoryError].GetTranslation(input.Language...)
			args["error"] = createErr.Error()
			fmt.Println(values.NewStringTemplate(text, args).Parse())
			return false
		}

		return checkInitProject(input)
	}
	if !f.IsDir() {
		_, text := i18nCollections[i18nInitProjectNotDirectory].GetTranslation(input.Language...)
		fmt.Println(values.NewStringTemplate(text, args).Parse())
		return false
	}

	// check if the project is already initialized
	modFile := filepath.Join(directory, "go.mod")
	if _, err := os.Stat(modFile); err == nil {
		_, text := i18nCollections[i18nInitExistProjectNotice].GetTranslation(input.Language...)
		fmt.Println(values.NewStringTemplate(text, args).Parse())
		LoadProject(input)
		return false
	}

	// try to write a file to the directory to check if it is writable
	_, tryCreate := os.Create(filepath.Join(directory, ".testing"))
	if tryCreate != nil {
		_, text := i18nCollections[i18nInitProjectCannotWrite].GetTranslation(input.Language...)
		fmt.Println(values.NewStringTemplate(text, args).Parse())
		return false
	}
	tryRemove := os.Remove(filepath.Join(directory, ".testing"))
	if tryRemove != nil {
		_, text := i18nCollections[i18nInitProjectCannotWrite].GetTranslation(input.Language...)
		fmt.Println(values.NewStringTemplate(text, args).Parse())
		return false
	}

	// write the input options to the configuration variable
	configuration.ProjectDirectory = directory
	configuration.PackageName = packageName
	return true
}

func printWorkingDirectory() {
	dir := configuration.ProjectDirectory
	if dir == "" {
		dir, _ = os.Getwd()
	}
	_, text := i18nCollections[i18nPrintWorkingDirectory].GetTranslation(cli.GetLanguage(cliOptions.PreferredLanguage, cliOptions.LanguageMapping)...)
	fmt.Println(values.NewStringTemplate(text, map[string]string{"directory": dir}).Parse())
}
