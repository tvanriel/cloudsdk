//go:build linux
// +build linux

package downwardsapi

type Configuration struct {
	// The root directory of the Kubernetes Downwards API mounts.
	Root string `hc:"root"`
}
