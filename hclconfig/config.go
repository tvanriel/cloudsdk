package hclconfig

import (
	"errors"
	"fmt"
	"os"
	"path"
)

var ErrFindConfigFile = errors.New("locate config file")

func HclConfiguration[T any](config T, projectName string) error {
	filename, ok := getFilename(projectName)
	if !ok {
		return ErrFindConfigFile
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read config file %q: %w", filename, err)
	}

	err = parseHCL(config, string(contents), filename)
	if err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}

func projectNameExpander(projectName string) func(string) string {
	return func(s string) string {
		if s == "PROJECT_NAME" {
			return projectName
		}
		return ""
	}
}

// DefaultConfigLocations makes a list of places where to search for the config file.
func DefaultConfigLocations(projectName string) []string {
	files := []string{
		"config.hcl",
		"$PROJECT_NAME.hcl",
	}

	paths := []string{
		".",
		"/etc/$PROJECT_NAME",
	}

	if homedir, err := os.UserHomeDir(); err == nil {
		paths = append(
			paths,
			path.Join(homedir, ".config/"),
			homedir,
		)
	}

	if xdgConfigDir := os.Getenv("XDG_CONFIG_DIR"); xdgConfigDir != "" {
		paths = append(paths, xdgConfigDir)
	}

	locations := make([]string, 0, len(files)*len(paths))
	for _, p := range paths {
		for _, f := range files {
			locations = append(locations, os.Expand(path.Join(p, f), projectNameExpander(projectName)))
		}
	}

	inResult := make(map[string]bool, len(locations))
	result := []string{}
	for _, l := range locations {
		if _, ok := inResult[l]; !ok {
			inResult[l] = true
			result = append(result, l)
		}
	}

	return result
}

func getFilename(projectName string) (string, bool) {
	tryFilenames := DefaultConfigLocations(projectName)

	for f := range tryFilenames {
		if _, err := os.Stat(tryFilenames[f]); err == nil {
			return tryFilenames[f], true
		}
	}
	return "", false
}
