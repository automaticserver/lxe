//go:build tools

package cli

import (
	// Importing a golangci lint package to track its version in go mod
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
