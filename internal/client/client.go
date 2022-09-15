package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Client is a wrapper for the Ockam binary.
type Client struct {
	path string
}

// NewClient returns a client which can be used to execute against the Ockam binary.
func NewClient() (*Client, error) {
	path, err := downloadBinary()
	if err != nil {
		return nil, fmt.Errorf("could not download ockam binary: %s", err)
	}

	c := &Client{path: path}
	return c, nil
}

// Run executes a command to the Ockam binary and returns the stdout results.
func (c *Client) Run(command ...string) (string, error) {
	cmd := exec.Command(c.path, command...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	return strings.TrimSuffix(out.String(), "\n"), err
}

// downloadBinary downloads the Ockam binary and places it in the user's cache.
func downloadBinary() (string, error) {
	var goarch string
	var goos string

	switch runtime.GOARCH {
	case amd64:
		goarch = x86_64
	case arm64:
		goarch = aarch64
	default:
		return "", fmt.Errorf("Arch `%s` is not supported by this provider", runtime.GOARCH)
	}
	switch runtime.GOOS {
	case darwin:
		goos = "apple-darwin"
	case linux:
		goos = "unknown-linux-gnu"
	default:
		return "", fmt.Errorf("OS `%s` is not supported by this provider", runtime.GOOS)

	}
	downloadURL := fmt.Sprintf(baseURL, version, goarch, goos)

	p, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	p = filepath.Join(p, "terraform-provider-ockam")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	p = filepath.Join(p, binary)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		u, err := url.Parse(downloadURL)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return "", err
		}

		file, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return "", err
		}
		defer file.Close()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()
		if _, err := io.Copy(file, resp.Body); err != nil {
			return "", err
		}
	}

	return p, nil
}
