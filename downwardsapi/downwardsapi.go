// The DownwardsAPI from Kubernetes is only available in Linux-based containers.
//go:build linux
// +build linux

package downwardsapi

import (
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

func NewDownwardsAPI(config Configuration) *DownwardsAPI {
	return &DownwardsAPI{
		root: config.Root,
	}
}

type DownwardsAPI struct {
	root string
}

type Field struct {
	Key   string
	Value string
}

func (e Field) String() string {
	return strings.Join([]string{
		e.Key,
		"=",
		strconv.Quote(e.Value),
	}, "")
}

func (d *DownwardsAPI) File(key string) ([]Field, error) {
	f, err := os.Open(path.Join(d.root, key))
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return Parse(string(b))
}
