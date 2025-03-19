// +build tools

package main

// This file lists tool dependencies for go modules.
// It ensures they're properly tracked in go.mod without being included in the build.
import (
	_ "github.com/axw/gocov/gocov"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/matm/gocov-html/cmd/gocov-html"
	_ "golang.org/x/tools/cmd/goimports"
)