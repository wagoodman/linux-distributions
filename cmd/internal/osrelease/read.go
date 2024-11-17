package osrelease

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read(file string) ([]Info, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	var osReleaseInfos []Info
	dec := json.NewDecoder(fh)
	if err := dec.Decode(&osReleaseInfos); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return osReleaseInfos, fh.Close()
}
