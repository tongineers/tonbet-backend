//go:build tools
// +build tools

package tools

import (
	_ "github.com/google/wire/cmd/wire"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
