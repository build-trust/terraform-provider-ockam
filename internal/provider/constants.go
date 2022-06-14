package provider

// Note This file if touched must be updated in ./github/constants.go.template
// for release automation to work

const (
	version = "0.61.0"
	baseURL = "https://github.com/ockam-network/ockam/releases/download/ockam_v%s/ockam.%s-%s"
	binary  = "ockam-v" + version
	arm64   = "arm64"
	amd64   = "amd64"
	aarch64 = "aarch64"
	x86_64  = "x86_64"
	darwin  = "darwin"
	linux   = "linux"
)
