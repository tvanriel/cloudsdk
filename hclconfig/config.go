package hclconfig

import (
	"errors"
	"os"
	"path"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func HclConfiguration[T any](config T, projectName string) error {
	filename, ok := getFilename(projectName)
	if !ok {
		return errors.New("error: cannot locate config file")
	}
	ctxVars, err := parseVariables(filename)
	err = hclsimple.DecodeFile(filename, &hcl.EvalContext{
		Variables: ctxVars,
	}, config)

	return err
}

func getFilename(projectName string) (string, bool) {
	tryFilenames := []string{"./config.hcl", projectName + ".hcl", "/etc/" + projectName + "/config.hcl"}
	if homedir, err := os.UserHomeDir(); err == nil {
		tryFilenames = append(
			tryFilenames,
			path.Join(homedir, ".config/"+projectName+".hcl"),
			path.Join(homedir, "."+projectName+".hcl"),
		)
	}
	if env := os.Getenv("XDG_CONFIG_DIR"); env != "" {
		tryFilenames = append(tryFilenames, env)
	}

	for f := range tryFilenames {
		if _, err := os.Open(tryFilenames[f]); err == nil {
			return tryFilenames[f], true
		}
	}
	return "", false
}
