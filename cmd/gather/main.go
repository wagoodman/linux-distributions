package main

import (
	"fmt"
	"github.com/wagoodman/linux-distributions/cmd/internal/eol"
	"github.com/wagoodman/linux-distributions/cmd/internal/osrelease"
	"path/filepath"
)

const (
	dataDir       = "data"
	osReleaseJSON = "os-release.json"
	eolJSON       = "eol.json"
)

func main() {
	ids, err := osrelease.Fetch(filepath.Join(dataDir, osReleaseJSON))
	if err != nil {
		panic(fmt.Errorf("failed to fetch os-release data: %w", err))
	}

	err = eol.Fetch(filepath.Join(dataDir, eolJSON), ids)
	if err != nil {
		panic(fmt.Errorf("failed to fetch eol data: %w", err))
	}
}
